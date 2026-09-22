package plugin

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
)

func init() { Register(videoPlugin()) }

//go:embed video/*.tmpl
var videoTemplates embed.FS

// videoPlugin converts uploaded clips into video every client can play.
//
// Design:
//
//   - ffmpeg does the work, in its own process, with arguments passed directly
//     and never through a shell. The output is one H.264 MP4 with its index at
//     the front, capped at 720 on the short side, and a JPEG poster: every
//     browser, iOS and Android plays that, so there is no player to pick per
//     platform and no HLS to serve.
//
//   - The videos table is the queue. A worker claims a pending row with a
//     conditional update, so a conversion survives a restart, replicas share
//     the work, and it runs the same with or without Redis. One conversion at
//     a time per replica, because ffmpeg uses every core it is given.
//
//   - Only real video containers are converted, read with a forced demuxer.
//     ffmpeg otherwise follows what a file says it is, and a playlist posing as
//     a video can name other files on the server for it to read.
//
//   - A video is created from an upload the caller owns, by its storage key,
//     so the key of someone else's file answers as not found.
//
// Found missing building the Instagram blueprint: the image chain covered
// photos, and there was nothing for a clip.
func videoPlugin() Plugin {
	return Plugin{
		Name:    "video",
		Version: "1.0.0",
		Summary: "Video uploads converted by ffmpeg to MP4 and a poster, for web and the Expo app",
		Description: `Turns uploaded clips into video every client can play.

  • POST /videos {key} converts an upload you own in the background; GET
    /videos/:id reports pending, processing, ready or failed; DELETE removes it
  • ffmpeg makes one H.264 MP4 (720 on the short side, playable before it has
    downloaded) and a poster frame, and records the size and length
  • The videos table is the queue: survives restarts, shared by replicas, no
    Redis needed; video.ready and video.failed reach the owner over realtime
  • Only MP4, MOV and WebM are converted, and a playlist posing as a video is
    refused, so a crafted file cannot make ffmpeg read files on the server
  • VideoPlayer for the web app and the Expo app (expo-video), with useVideo
    polling until the conversion is done

Needs ffmpeg: the API's Docker image gets it; install it locally to convert in
development (winget install ffmpeg, brew install ffmpeg, apt install ffmpeg).`,

		NextSteps: []string{
			"Run the migration:   grit migrate",
			"Install packages:    pnpm install",
			"Install ffmpeg to convert locally (the Docker image already has it): winget install ffmpeg | brew install ffmpeg | apt install ffmpeg",
			"Upload with accepts=video, then POST /api/v1/videos {\"key\": \"<the upload's key>\"}",
			"Play it: <VideoPlayer id={video.id} /> from @/components/video-player, on web and in the Expo app",
			"Limits: video.NewService (in routes.go) sets 720p and three minutes; change videoService.Options.MaxEdge and MaxDuration after it",
		},

		NodeDeps: []Dependency{
			// The version Expo SDK 54 bundles (expo/bundledNativeModules.json).
			{Name: "expo-video", Version: "~3.0.16", Workspace: "apps/expo"},
		},

		Files:      videoFiles,
		Injections: videoInjections,
	}
}

func videoTemplate(ctx Context, name string) string {
	raw, err := videoTemplates.ReadFile("video/" + name)
	if err != nil {
		// The templates are compiled in; a missing one is a build mistake.
		panic("video plugin template missing: " + name)
	}
	return strings.ReplaceAll(string(raw), "{{MODULE}}", ctx.Module)
}

func hasWebApp(ctx Context) bool {
	info, err := os.Stat(filepath.Join(ctx.Root, "apps", "web"))
	return err == nil && info.IsDir()
}

func videoFiles(ctx Context) map[string]string {
	files := map[string]string{
		apiPath(ctx, "internal/models/video.go"):     videoTemplate(ctx, "video_model.go.tmpl"),
		apiPath(ctx, "internal/video/ffmpeg.go"):     videoTemplate(ctx, "ffmpeg.go.tmpl"),
		apiPath(ctx, "internal/video/service.go"):    videoTemplate(ctx, "video_service.go.tmpl"),
		apiPath(ctx, "internal/video/video_test.go"): videoTemplate(ctx, "video_test.go.tmpl"),
		apiPath(ctx, "internal/handlers/video.go"):   videoTemplate(ctx, "video_handler.go.tmpl"),
	}
	if hasWebApp(ctx) {
		files["apps/web/lib/video.ts"] = videoTemplate(ctx, "web_video.ts.tmpl")
		files["apps/web/hooks/use-video.ts"] = videoTemplate(ctx, "web_use_video.ts.tmpl")
		files["apps/web/components/video-player.tsx"] = videoTemplate(ctx, "web_video_player.tsx.tmpl")
	}
	if hasExpo(ctx) {
		files["apps/expo/lib/video.ts"] = videoTemplate(ctx, "expo_video.ts.tmpl")
		files["apps/expo/hooks/use-video.ts"] = videoTemplate(ctx, "expo_use_video.ts.tmpl")
		files["apps/expo/components/video-player.tsx"] = videoTemplate(ctx, "expo_video_player.tsx.tmpl")
	}
	return files
}

func videoInjections(ctx Context) []Injection {
	return []Injection{
		{
			File:   apiPath(ctx, "internal/models/user.go"),
			Marker: "// grit:models",
			Code:   "\t\t&Video{},",
		},
		{
			File:   apiPath(ctx, "internal/routes/routes.go"),
			Marker: "// grit:imports",
			Code:   "\t\"" + ctx.Module + "/internal/video\"",
		},
		{
			File:   apiPath(ctx, "internal/routes/routes.go"),
			Marker: "// grit:handlers",
			Code: "\t// Videos, converted by ffmpeg in the background; the owner hears when\n" +
				"\t// one is ready or failed on their realtime channel.\n" +
				"\tvideoHandler := handlers.NewVideoHandler(nil)\n" +
				"\tif svc.Storage != nil {\n" +
				"\t\tvideoService := video.NewService(db, svc.Storage.Disk())\n" +
				"\t\tvideoService.OnReady = func(v models.Video) {\n" +
				"\t\t\trealtimeHub.SendToUsers([]string{v.UserID}, realtime.Event{Type: \"video.ready\", Payload: v})\n" +
				"\t\t}\n" +
				"\t\tvideoService.OnFailed = func(v models.Video) {\n" +
				"\t\t\trealtimeHub.SendToUsers([]string{v.UserID}, realtime.Event{Type: \"video.failed\", Payload: v})\n" +
				"\t\t}\n" +
				"\t\tvideoService.Start(context.Background())\n" +
				"\t\tvideoHandler.Videos = videoService\n" +
				"\t}",
		},
		{
			File:   apiPath(ctx, "internal/routes/routes.go"),
			Marker: "// grit:routes:protected",
			Code: "\t\t// Videos: convert an upload you own, then poll or listen for video.ready.\n" +
				"\t\tprotected.POST(\"/videos\", videoHandler.Create)\n" +
				"\t\tprotected.GET(\"/videos/:id\", videoHandler.Get)\n" +
				"\t\tprotected.DELETE(\"/videos/:id\", videoHandler.Delete)",
		},
		{
			// The runtime image gets ffmpeg. A Dockerfile from before the
			// marker is reported, with the line to add, rather than failing.
			File:     apiPath(ctx, "Dockerfile"),
			Marker:   "# grit:runtime-packages",
			Code:     "RUN apk --no-cache add ffmpeg",
			Optional: true,
		},
	}
}
