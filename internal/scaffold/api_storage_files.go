package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeStorageFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "storage", "storage.go"):      storageServiceGo(),
		filepath.Join(apiRoot, "internal", "storage", "disk.go"):         storageDiskGo(),
		filepath.Join(apiRoot, "internal", "storage", "local.go"):        storageLocalGo(),
		filepath.Join(apiRoot, "internal", "storage", "disk_test.go"):    storageDiskTestGo(),
		filepath.Join(apiRoot, "internal", "storage", "store.go"):        storageStoreGo(),
		filepath.Join(apiRoot, "internal", "storage", "store_test.go"):   storageStoreTestGo(),
		filepath.Join(apiRoot, "internal", "storage", "disks.go"):        storageDisksGo(),
		filepath.Join(apiRoot, "internal", "storage", "image.go"):        storageImageGo(),
		filepath.Join(apiRoot, "internal", "storage", "image_test.go"):   storageImageTestGo(),
		filepath.Join(apiRoot, "internal", "storage", "url_test.go"):     storageURLTestGo(module),
		filepath.Join(apiRoot, "internal", "handlers", "upload.go"):      uploadHandlerGo(),
		filepath.Join(apiRoot, "internal", "services", "upload.go"):      apiUploadServiceGo(),
		filepath.Join(apiRoot, "internal", "files", "file_ref.go"):       filesFileRefGo(),
		filepath.Join(apiRoot, "internal", "files", "accepts.go"):        filesAcceptsGo(),
		filepath.Join(apiRoot, "internal", "files", "file_ref_test.go"):  filesFileRefTestGo(),
		filepath.Join(apiRoot, "internal", "files", "lifecycle.go"):      filesLifecycleGo(),
		filepath.Join(apiRoot, "internal", "files", "lifecycle_test.go"): filesLifecycleTestGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func storageServiceGo() string {
	return `package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"{{MODULE}}/internal/config"
)

// Storage is the file store the app started with.
//
// It holds one Disk, chosen by STORAGE_DRIVER: an S3-compatible bucket (minio,
// s3, r2, b2) or a directory on this machine (local). New code can take the
// Disk itself from Disk(). The methods below keep the names handlers have
// always called, each a thin wrapper over the Disk, so code written before
// the interface existed compiles and behaves as it did.
type Storage struct {
	disk Disk
}

// PublicPrefixes are the key prefixes anyone may read without a signature:
// uploaded files and their thumbnails, which pages link to directly. Backups,
// private originals and every other key are read through GetSignedURL.
//
// STORAGE_PUBLIC_PREFIXES replaces them, through SetPublicPrefixes.
var PublicPrefixes = []string{"uploads/", "thumbnails/"}

// SetPublicPrefixes replaces PublicPrefixes. Call it before opening a store: a
// bucket's policy is written from the prefixes when it connects. Each prefix
// is given a trailing slash, so "media" cannot make "media-private/" public,
// and an empty list keeps the prefixes there are.
func SetPublicPrefixes(prefixes []string) {
	cleaned := make([]string, 0, len(prefixes))
	for _, prefix := range prefixes {
		if prefix = strings.Trim(strings.TrimSpace(prefix), "/"); prefix != "" {
			cleaned = append(cleaned, prefix+"/")
		}
	}
	if len(cleaned) > 0 {
		PublicPrefixes = cleaned
	}
}

// IsPublicKey reports whether key is under one of PublicPrefixes.
func IsPublicKey(key string) bool {
	for _, prefix := range PublicPrefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

// BucketPolicy allows anonymous reads under PublicPrefixes, and nothing else.
func BucketPolicy(bucket string) string {
	resources := make([]string, 0, len(PublicPrefixes))
	for _, prefix := range PublicPrefixes {
		resources = append(resources, strconv.Quote("arn:aws:s3:::"+bucket+"/"+prefix+"*"))
	}
	return "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"AWS\":[\"*\"]}," +
		"\"Action\":[\"s3:GetObject\"],\"Resource\":[" + strings.Join(resources, ",") + "]}]}"
}

// New connects to an S3-compatible bucket: AWS S3, MinIO, Cloudflare R2 or
// Backblaze B2.
func New(cfg config.StorageConfig) (*Storage, error) {
	disk, err := NewS3Disk(cfg)
	if err != nil {
		return nil, err
	}
	return &Storage{disk: disk}, nil
}

// NewLocal keeps files in a directory on this machine (STORAGE_DRIVER=local).
func NewLocal(cfg LocalConfig) (*Storage, error) {
	disk, err := NewLocalDisk(cfg)
	if err != nil {
		return nil, err
	}
	return &Storage{disk: disk}, nil
}

// Wrap returns a Storage over any Disk: a driver of your own, or a fake in a
// test.
func Wrap(disk Disk) *Storage {
	return &Storage{disk: disk}
}

// Disk is the driver behind this store.
func (s *Storage) Disk() Disk {
	return s.disk
}

// Describe says where this store keeps files, for a startup log line.
func (s *Storage) Describe() string {
	switch d := s.disk.(type) {
	case *LocalDisk:
		return fmt.Sprintf("local disk at %s, served from %s", d.root, d.publicURL)
	case *S3Disk:
		if d.cfg.Endpoint == "" {
			return fmt.Sprintf("bucket %q on AWS S3 (%s)", d.bucket, d.cfg.Region)
		}
		return fmt.Sprintf("bucket %q at %s", d.bucket, d.cfg.Endpoint)
	default:
		return fmt.Sprintf("%T", d)
	}
}

// Driver names the kind of store this is: local, s3, or the type behind a Disk
// of your own.
//
// For /api/health, which is public so a load balancer can reach it. Describe
// also names the bucket and the endpoint, which is right for a startup log line
// and wrong for an unauthenticated response.
func (s *Storage) Driver() string {
	switch d := s.disk.(type) {
	case *LocalDisk:
		return "local"
	case *S3Disk:
		return "s3"
	default:
		return fmt.Sprintf("%T", d)
	}
}

// Upload stores a file at the given key.
func (s *Storage) Upload(ctx context.Context, key string, reader io.Reader, contentType string) error {
	return s.disk.Put(ctx, key, reader, PutOptions{ContentType: contentType})
}

// Download opens a stored file. The caller closes it.
func (s *Storage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	return s.disk.Get(ctx, key)
}

// Delete removes a stored file.
func (s *Storage) Delete(ctx context.Context, key string) error {
	return s.disk.Delete(ctx, key)
}

// DeleteMany removes many stored files. A key that is already gone is not an
// error.
func (s *Storage) DeleteMany(ctx context.Context, keys []string) error {
	return s.disk.Delete(ctx, keys...)
}

// GetURL returns the URL a browser loads a public file from.
func (s *Storage) GetURL(key string) string {
	return s.disk.URL(key)
}

// GetSignedURL returns a link to any file that stops working after duration.
func (s *Storage) GetSignedURL(ctx context.Context, key string, duration time.Duration) (string, error) {
	return s.disk.TemporaryURL(ctx, key, duration)
}

// Stat returns the size and content type of a stored file.
func (s *Storage) Stat(ctx context.Context, key string) (int64, string, error) {
	obj, err := s.disk.Stat(ctx, key)
	if err != nil {
		return 0, "", err
	}
	return obj.Size, obj.ContentType, nil
}

// PresignPutURL generates a pre-signed PUT URL for a direct browser upload,
// valid for an hour. It returns ErrPresignUnsupported when the driver takes
// uploads through the API instead, and the upload handler tells the client
// to send the file to POST /uploads.
func (s *Storage) PresignPutURL(ctx context.Context, key, contentType string, contentLength int64) (string, error) {
	presigner, ok := s.disk.(Presigner)
	if !ok {
		return "", ErrPresignUnsupported
	}
	return presigner.PresignPut(ctx, key, contentType, contentLength, time.Hour)
}

// FileServer is the handler that serves this store's files from the API, or
// nil when something else serves them (a bucket serves its own).
func (s *Storage) FileServer() http.Handler {
	if handler, ok := s.disk.(http.Handler); ok {
		return handler
	}
	return nil
}

// S3Disk is a Disk over one S3-compatible bucket.
type S3Disk struct {
	client *s3.Client
	bucket string
	cfg    config.StorageConfig
}

// NewS3Disk connects to the bucket in cfg, creating it if it does not exist.
func NewS3Disk(cfg config.StorageConfig) (*S3Disk, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if cfg.Endpoint != "" {
				return aws.Endpoint{
					URL:               cfg.Endpoint,
					HostnameImmutable: true,
					SigningRegion:     cfg.Region,
				}, nil
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		},
	)

	loadOptions := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
	}
	// With no key, the SDK's own chain finds credentials: AWS_* variables, a
	// shared profile, or the IAM role of the EC2, ECS or Lambda it runs on.
	if cfg.AccessKey != "" {
		loadOptions = append(loadOptions, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		))
	}
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), loadOptions...)
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %w", err)
	}

	// Path-style addressing is required for MinIO and works for R2 / B2.
	// AWS S3 buckets created after Sep 2020 reject path-style and require
	// virtual-hosted style. We use the endpoint as the signal: an empty
	// endpoint means "go to default AWS regional endpoint" = real S3 =
	// virtual-hosted. A non-empty endpoint means a third-party S3-clone
	// that needs path-style.
	usePathStyle := cfg.Endpoint != ""
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = usePathStyle
	})

	// Verify bucket exists with a quick head request
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(cfg.Bucket),
	})
	if err != nil {
		// Try to create the bucket
		_, createErr := client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(cfg.Bucket),
		})
		if createErr != nil {
			return nil, fmt.Errorf("bucket %q not accessible and cannot be created: %w", cfg.Bucket, err)
		}
	}

	// Anyone may read what an <img> or a download link points at, and nothing
	// else. The policy used to cover every key, backups and private originals
	// included, so a backup's key seen once in a log or a Referer header was a
	// permanent anonymous download of the database. PutBucketPolicy replaces the
	// policy a bucket has, so an existing bucket is narrowed on the next start.
	if _, err := client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
		Bucket: aws.String(cfg.Bucket),
		Policy: aws.String(BucketPolicy(cfg.Bucket)),
	}); err != nil {
		// Cloudflare R2 and Backblaze B2 have no bucket policies: public access
		// is switched on per bucket in their dashboards. Anywhere else, a refusal
		// leaves the bucket with whatever policy it had before.
		log.Printf("storage: could not set the bucket policy on %q: %v. Where the provider has bucket policies, allow anonymous reads on %s only; "+
			"where it has none (R2, B2), every key in a public bucket is public, so put backups in a private bucket of their own (STORAGE_DISKS=backups) and serve private files with a temporary URL",
			cfg.Bucket, err, strings.Join(PublicPrefixes, ", "))
	}

	return &S3Disk{
		client: client,
		bucket: cfg.Bucket,
		cfg:    cfg,
	}, nil
}

// s3Failed wraps an SDK error, reporting a missing object as ErrNotFound.
func s3Failed(op, key string, err error) error {
	var noSuchKey *types.NoSuchKey
	var notFound *types.NotFound
	var response *awshttp.ResponseError
	if errors.As(err, &noSuchKey) || errors.As(err, &notFound) ||
		(errors.As(err, &response) && response.HTTPStatusCode() == http.StatusNotFound) {
		return fmt.Errorf("%s %q: %w (%w)", op, key, ErrNotFound, err)
	}
	return fmt.Errorf("%s %q: %w", op, key, err)
}

// Put stores r in the bucket at key.
func (d *S3Disk) Put(ctx context.Context, key string, r io.Reader, opts PutOptions) error {
	if err := checkKey(key); err != nil {
		return err
	}
	if err := checkVisibility(key, opts.Visibility); err != nil {
		return err
	}
	contentType := opts.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	_, err := d.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(d.bucket),
		Key:         aws.String(key),
		Body:        r,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("uploading %q: %w", key, err)
	}
	return nil
}

// Get opens the object at key.
func (d *S3Disk) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	if err := checkKey(key); err != nil {
		return nil, err
	}
	result, err := d.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, s3Failed("downloading", key, err)
	}
	return result.Body, nil
}

// GetRange opens length bytes of the object at key from offset, so ServeFile
// answers a Range request with one ranged GET rather than the whole object.
func (d *S3Disk) GetRange(ctx context.Context, key string, offset, length int64) (io.ReadCloser, error) {
	if err := checkKey(key); err != nil {
		return nil, err
	}
	if offset < 0 || length <= 0 {
		return nil, fmt.Errorf("storage: invalid range %d+%d for %q", offset, length, key)
	}
	result, err := d.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(key),
		Range:  aws.String(fmt.Sprintf("bytes=%d-%d", offset, offset+length-1)),
	})
	if err != nil {
		return nil, s3Failed("downloading", key, err)
	}
	return result.Body, nil
}

// Exists reports whether an object is stored at key.
func (d *S3Disk) Exists(ctx context.Context, key string) (bool, error) {
	return exists(ctx, d, key)
}

// Stat asks the bucket what it holds at key.
//
// Needed because a presigned upload never passes through this server: the only
// way to know what actually landed in the bucket is to ask the bucket.
func (d *S3Disk) Stat(ctx context.Context, key string) (Object, error) {
	if err := checkKey(key); err != nil {
		return Object{}, err
	}
	out, err := d.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return Object{}, s3Failed("stat", key, err)
	}
	return Object{
		Key:          key,
		Size:         aws.ToInt64(out.ContentLength),
		ContentType:  aws.ToString(out.ContentType),
		LastModified: aws.ToTime(out.LastModified),
	}, nil
}

// Delete removes keys from the bucket, 1,000 to a request, the most one
// DeleteObjects call takes. A key that is already gone is not an error.
func (d *S3Disk) Delete(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		if err := checkKey(key); err != nil {
			return err
		}
	}
	if len(keys) == 1 {
		_, err := d.client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(d.bucket),
			Key:    aws.String(keys[0]),
		})
		if err != nil {
			return fmt.Errorf("deleting %q: %w", keys[0], err)
		}
		return nil
	}
	for start := 0; start < len(keys); start += 1000 {
		end := start + 1000
		if end > len(keys) {
			end = len(keys)
		}
		objects := make([]types.ObjectIdentifier, 0, end-start)
		for _, key := range keys[start:end] {
			objects = append(objects, types.ObjectIdentifier{Key: aws.String(key)})
		}
		out, err := d.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(d.bucket),
			Delete: &types.Delete{Objects: objects, Quiet: aws.Bool(true)},
		})
		if err != nil {
			return fmt.Errorf("deleting %d objects: %w", len(objects), err)
		}
		if len(out.Errors) > 0 {
			first := out.Errors[0]
			return fmt.Errorf("deleting %d of %d objects failed, the first %q: %s",
				len(out.Errors), len(objects), aws.ToString(first.Key), aws.ToString(first.Message))
		}
	}
	return nil
}

// Copy copies the object at from to to, inside the bucket.
func (d *S3Disk) Copy(ctx context.Context, from, to string) error {
	if err := checkKey(from); err != nil {
		return err
	}
	if err := checkKey(to); err != nil {
		return err
	}
	_, err := d.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(d.bucket),
		CopySource: aws.String(d.bucket + "/" + escapeKey(from)),
		Key:        aws.String(to),
	})
	if err != nil {
		return s3Failed("copying", from, err)
	}
	return nil
}

// Move copies the object at from to to, then deletes from. A bucket has no
// rename.
func (d *S3Disk) Move(ctx context.Context, from, to string) error {
	if err := d.Copy(ctx, from, to); err != nil {
		return err
	}
	return d.Delete(ctx, from)
}

// List returns every object whose key starts with prefix, in key order.
func (d *S3Disk) List(ctx context.Context, prefix string) ([]Object, error) {
	var objects []Object
	pages := s3.NewListObjectsV2Paginator(d.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(d.bucket),
		Prefix: aws.String(prefix),
	})
	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("listing %q: %w", prefix, err)
		}
		for _, item := range page.Contents {
			objects = append(objects, Object{
				Key:          aws.ToString(item.Key),
				Size:         aws.ToInt64(item.Size),
				LastModified: aws.ToTime(item.LastModified),
			})
		}
	}
	return objects, nil
}

// URL returns the URL a browser should load this object from.
//
// The SDK endpoint and the browser-facing origin are not always the same host.
// MinIO serves objects from the host it takes API calls on, so the default
// (<endpoint>/<bucket>/<key>) is right there. Cloudflare R2 is the case that
// breaks: <account>.r2.cloudflarestorage.com only answers SigV4-signed
// requests, so an <img> pointed at it gets a 401: the upload succeeds and
// nothing ever renders, which reads like a CORS problem and is not one.
//
// Setting R2_PUBLIC_URL (or STORAGE_PUBLIC_URL) switches this to
// <PublicURL>/<key>. Public origins (an r2.dev subdomain, a custom domain, a
// CDN in front of S3) are already scoped to one bucket, so the bucket
// segment is deliberately not repeated.
func (d *S3Disk) URL(key string) string {
	escaped := escapeKey(key)
	if public := strings.TrimRight(d.cfg.PublicURL, "/"); public != "" {
		return fmt.Sprintf("%s/%s", public, escaped)
	}
	endpoint := strings.TrimRight(d.cfg.Endpoint, "/")
	if endpoint == "" {
		// AWS S3 at its regional default, where a bucket is addressed by host.
		region := d.cfg.Region
		if region == "" {
			region = "us-east-1"
		}
		return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", d.bucket, region, escaped)
	}
	return fmt.Sprintf("%s/%s/%s", endpoint, d.bucket, escaped)
}

// TemporaryURL returns a pre-signed GET URL valid for ttl.
func (d *S3Disk) TemporaryURL(ctx context.Context, key string, ttl time.Duration) (string, error) {
	if err := checkKey(key); err != nil {
		return "", err
	}
	presigner := s3.NewPresignClient(d.client)
	result, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", fmt.Errorf("generating signed URL for %q: %w", key, err)
	}
	return result.URL, nil
}

// PresignPut generates a pre-signed PUT URL for a direct browser upload.
//
// size is signed into the URL, so S3 rejects a PUT of any other size.
// Without it the URL is an unbounded write capability: a client can ask to
// upload two megabytes and then send five gigabytes, and nothing on this side
// ever sees it happen.
//
// It is an exact match rather than a ceiling, which the client can satisfy
// because it optimises the image first and therefore knows the byte count
// before it asks for a URL.
func (d *S3Disk) PresignPut(ctx context.Context, key, contentType string, size int64, ttl time.Duration) (string, error) {
	if err := checkKey(key); err != nil {
		return "", err
	}
	presigner := s3.NewPresignClient(d.client)
	result, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		ContentLength: aws.Int64(size),
		Bucket:        aws.String(d.bucket),
		Key:           aws.String(key),
		ContentType:   aws.String(contentType),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", fmt.Errorf("generating presigned PUT URL for %q: %w", key, err)
	}
	return result.URL, nil
}
`
}

func storageImageGo() string {
	return `package storage

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"

	"{{MODULE}}/internal/media"
)

// These helpers are the path an image takes when the upload pipeline did not
// handle it: a presigned upload, which goes straight to storage, or a type the
// pipeline skipped. They run through media.Transform, the pipeline itself, so
// EXIF orientation is applied and metadata stripped the same way on both paths.
// Before, they decoded without orientation, and a portrait phone photo got a
// sideways thumbnail.

// MaxImageWidth is the maximum width for processed images.
const MaxImageWidth = 1920

// ThumbnailSize is the edge of the square thumbnail GenerateThumbnail makes:
// the "thumb" rendition of media's default profile, so a thumbnail is the same
// size whichever path made it.
const ThumbnailSize = 400

// MaxImageBytes is the most the helpers read of an image. Uploads are capped at
// 50 MB, and the one byte over tells a larger file from one exactly at the cap.
const MaxImageBytes = 50<<20 + 1

var (
	// ErrImageTooLarge is returned for an image refused before decoding: over
	// MaxImageBytes, or with a header claiming more pixels than the media
	// profile allows. A few kilobytes of PNG can claim 30000x30000, which
	// decodes to 3.6 GB and takes the process with it.
	ErrImageTooLarge = errors.New("image is too large to decode")
	// ErrUnreadableImage is returned for data that is not an image this
	// package can decode. Like ErrImageTooLarge it fails the same way every
	// time, so a job should not retry it.
	ErrUnreadableImage = errors.New("image cannot be decoded")
)

// readImage reads an image and checks its header, refusing it when the header
// claims more pixels than media.Get("").MaxPixels. Decoding commits the memory
// for every pixel before it reads them, so the header is the only place a
// decompression bomb can be stopped, and it is stopped here with an error a job
// knows not to retry.
func readImage(reader io.Reader) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(reader, MaxImageBytes))
	if err != nil {
		return nil, fmt.Errorf("reading image: %w", err)
	}
	if len(data) >= MaxImageBytes {
		return nil, fmt.Errorf("image is over %d MB: %w", (MaxImageBytes-1)>>20, ErrImageTooLarge)
	}
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("reading image header: %v: %w", err, ErrUnreadableImage)
	}
	limit := int64(media.Get("").MaxPixels)
	if px := int64(cfg.Width) * int64(cfg.Height); limit > 0 && px > limit {
		return nil, fmt.Errorf("image is %dx%d, over the %d megapixel limit: %w",
			cfg.Width, cfg.Height, limit/1000000, ErrImageTooLarge)
	}
	return data, nil
}

// transformImage runs one image through media.Transform at size, keeping the
// format it came in: PNG stays PNG, anything else becomes JPEG.
func transformImage(reader io.Reader, mimeType string, size media.Size) ([]byte, error) {
	data, err := readImage(reader)
	if err != nil {
		return nil, err
	}
	profile := media.Get("")
	profile.Max = size
	profile.Quality = 0.85
	profile.Format = media.JPEG
	if strings.EqualFold(mimeType, "image/png") {
		profile.Format = media.PNG
	}
	// No extra renditions: the caller wants this one image.
	profile.Renditions = map[string]media.Size{}
	result, err := media.Transform(bytes.NewReader(data), profile)
	if err != nil {
		return nil, fmt.Errorf("processing image: %v: %w", err, ErrUnreadableImage)
	}
	return result.Primary.Bytes, nil
}

// ProcessImage resizes an image wider than MaxImageWidth, keeping its aspect
// ratio and applying its EXIF orientation, and returns the encoded bytes.
func ProcessImage(reader io.Reader, mimeType string) ([]byte, error) {
	// The bound is the width: the height limit is far past any image under the
	// pixel limit that is also this wide.
	return transformImage(reader, mimeType, media.Fit(MaxImageWidth, MaxImageWidth*64))
}

// GenerateThumbnail makes a ThumbnailSize square thumbnail, cropped from the
// centre of the image the right way up.
func GenerateThumbnail(reader io.Reader, mimeType string) ([]byte, error) {
	return transformImage(reader, mimeType, media.Fill(ThumbnailSize, ThumbnailSize))
}

// IsImageMimeType returns true if the MIME type is a supported image format.
func IsImageMimeType(mimeType string) bool {
	switch strings.ToLower(mimeType) {
	case "image/jpeg", "image/png", "image/gif":
		return true
	}
	return false
}
`
}

func uploadHandlerGo() string { return tmpl("api/handlers/upload.go") }

// storageURLTestGo guards the object-URL rule that makes images actually
// render.
//
// The failure it exists for is quiet: point object URLs at R2's S3 endpoint
// and every upload succeeds while every <img> gets a 401, which reads like a
// CORS problem and sends people to the wrong setting for an afternoon.
func storageURLTestGo(module string) string {
	return fmt.Sprintf(`package storage

import (
	"encoding/json"
	"testing"

	"%s/internal/config"
)

// Only uploaded files and their thumbnails are public. The policy used to cover
// every key, database backups included.
func TestBucketPolicyIsScopedToPublicPrefixes(t *testing.T) {
	var doc struct {
		Statement []struct {
			Resource []string
		}
	}
	if err := json.Unmarshal([]byte(BucketPolicy("b")), &doc); err != nil {
		t.Fatalf("the policy is not JSON: %%v", err)
	}
	if len(doc.Statement) != 1 {
		t.Fatalf("want one statement, got %%d", len(doc.Statement))
	}
	want := map[string]bool{"arn:aws:s3:::b/uploads/*": true, "arn:aws:s3:::b/thumbnails/*": true}
	for _, r := range doc.Statement[0].Resource {
		if !want[r] {
			t.Errorf("the policy makes %%s public", r)
		}
		delete(want, r)
	}
	if len(want) != 0 {
		t.Errorf("missing from the policy: %%v", want)
	}
}

func TestGetURL(t *testing.T) {
	cases := []struct {
		name string
		cfg  config.StorageConfig
		key  string
		want string
	}{
		{
			name: "no public origin keeps the bucket segment (MinIO)",
			cfg:  config.StorageConfig{Endpoint: "http://localhost:9002", Bucket: "uploads"},
			key:  "uploads/2026/08/a.png",
			want: "http://localhost:9002/uploads/uploads/2026/08/a.png",
		},
		{
			// R2's S3 endpoint only answers signed requests, so object URLs
			// must come from the bucket's public origin instead.
			name: "public origin replaces endpoint and drops the bucket",
			cfg: config.StorageConfig{
				Endpoint:  "https://acct.r2.cloudflarestorage.com",
				Bucket:    "uploads",
				PublicURL: "https://pub-abc123.r2.dev",
			},
			key:  "uploads/2026/08/a.png",
			want: "https://pub-abc123.r2.dev/uploads/2026/08/a.png",
		},
		{
			name: "trailing slash on the public origin is not doubled",
			cfg:  config.StorageConfig{PublicURL: "https://cdn.example.com/"},
			key:  "uploads/a.png",
			want: "https://cdn.example.com/uploads/a.png",
		},
		{
			name: "spaces in a key stay escaped",
			cfg:  config.StorageConfig{PublicURL: "https://cdn.example.com"},
			key:  "uploads/my file.png",
			want: "https://cdn.example.com/uploads/my%%20file.png",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := Wrap(&S3Disk{bucket: c.cfg.Bucket, cfg: c.cfg})
			if got := s.GetURL(c.key); got != c.want {
				t.Errorf("GetURL(%%q) = %%q, want %%q", c.key, got, c.want)
			}
		})
	}
}
`, module)
}

// storageImageTestGo is the regression test for H12 in the contact-app review,
// shipped so it runs where the code does.
func storageImageTestGo() string {
	return `package storage

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"runtime"
	"testing"

	"{{MODULE}}/internal/media"
)

// pngClaiming is a PNG whose header claims w x h RGBA pixels, with one short
// row of data behind it. Decoding it allocates for every claimed pixel before
// discovering the data is not there.
func pngClaiming(w, h uint32) []byte {
	var b bytes.Buffer
	b.WriteString("\x89PNG\r\n\x1a\n")
	chunk := func(kind string, data []byte) {
		_ = binary.Write(&b, binary.BigEndian, uint32(len(data)))
		b.WriteString(kind)
		b.Write(data)
		_ = binary.Write(&b, binary.BigEndian, crc32.ChecksumIEEE(append([]byte(kind), data...)))
	}
	ihdr := make([]byte, 13)
	binary.BigEndian.PutUint32(ihdr[0:], w)
	binary.BigEndian.PutUint32(ihdr[4:], h)
	ihdr[8], ihdr[9] = 8, 6
	chunk("IHDR", ihdr)
	var z bytes.Buffer
	zw := zlib.NewWriter(&z)
	_, _ = zw.Write(make([]byte, 64))
	_ = zw.Close()
	chunk("IDAT", z.Bytes())
	chunk("IEND", nil)
	return b.Bytes()
}

// A few kilobytes claiming 30000x30000 would allocate 3.6 GB in image.Decode.
// The thumbnail worker ran that in the API process and retried it five times.
func TestAPixelBombIsRefusedBeforeDecoding(t *testing.T) {
	bomb := pngClaiming(30000, 30000)
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	_, err := GenerateThumbnail(bytes.NewReader(bomb), "image/png")
	runtime.ReadMemStats(&after)
	if !errors.Is(err, ErrImageTooLarge) {
		t.Fatalf("want ErrImageTooLarge, got %v", err)
	}
	if grew := after.TotalAlloc - before.TotalAlloc; grew > 64<<20 {
		t.Errorf("refusing the image allocated %d MB", grew>>20)
	}
	if _, err := ProcessImage(bytes.NewReader(bomb), "image/png"); !errors.Is(err, ErrImageTooLarge) {
		t.Errorf("ProcessImage: want ErrImageTooLarge, got %v", err)
	}
}

func TestNotAnImageIsUnreadable(t *testing.T) {
	if _, err := GenerateThumbnail(bytes.NewReader([]byte("not an image")), "image/png"); !errors.Is(err, ErrUnreadableImage) {
		t.Errorf("want ErrUnreadableImage, got %v", err)
	}
}

func TestAnOrdinaryImageStillMakesAThumbnail(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 640, 480))
	for x := 0; x < 640; x++ {
		img.Set(x, x%480, color.RGBA{R: 255, A: 255})
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	thumb, err := GenerateThumbnail(bytes.NewReader(buf.Bytes()), "image/png")
	if err != nil {
		t.Fatal(err)
	}
	got, _, err := image.DecodeConfig(bytes.NewReader(thumb))
	if err != nil || got.Width != ThumbnailSize || got.Height != ThumbnailSize {
		t.Errorf("thumbnail is %dx%d (%v), want %dx%d", got.Width, got.Height, err, ThumbnailSize, ThumbnailSize)
	}
	if thumb := media.DefaultProfile().Renditions["thumb"]; thumb.Width != ThumbnailSize || thumb.Height != ThumbnailSize {
		t.Errorf("ThumbnailSize is %d, the media pipeline's thumb is %dx%d", ThumbnailSize, thumb.Width, thumb.Height)
	}
}

// jpegWithOrientation encodes img as a JPEG carrying an EXIF Orientation tag,
// as a phone writes a portrait photo: the pixels stay landscape and the tag
// says how to turn them.
func jpegWithOrientation(t *testing.T, img image.Image, orientation byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 95}); err != nil {
		t.Fatal(err)
	}
	src := buf.Bytes()
	// A big-endian TIFF header and one IFD entry: tag 0x0112, SHORT, count 1.
	tiff := []byte{'M', 'M', 0, 42, 0, 0, 0, 8, 0, 1, 0x01, 0x12, 0, 3, 0, 0, 0, 1, 0, orientation, 0, 0, 0, 0, 0, 0}
	payload := append([]byte("Exif\x00\x00"), tiff...)
	out := append([]byte{}, src[:2]...)
	out = append(out, 0xFF, 0xE1, byte((len(payload)+2)>>8), byte(len(payload)+2))
	out = append(out, payload...)
	return append(out, src[2:]...)
}

func isRed(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	return r > 0xa000 && g < 0x6000 && b < 0x6000
}

func isBlue(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	return b > 0xa000 && r < 0x6000 && g < 0x6000
}

// A photo whose EXIF says "turn me 90 degrees clockwise" gets a thumbnail the
// right way up. The helpers used to decode without orientation, so a presigned
// upload or a skipped one got a sideways thumbnail.
func TestThumbnailsFollowEXIFOrientation(t *testing.T) {
	// Landscape pixels: red on the left, blue on the right. Turned clockwise,
	// red is on top.
	wide := image.NewRGBA(image.Rect(0, 0, 160, 80))
	for y := 0; y < 80; y++ {
		for x := 0; x < 160; x++ {
			c := color.RGBA{R: 255, A: 255}
			if x >= 80 {
				c = color.RGBA{B: 255, A: 255}
			}
			wide.Set(x, y, c)
		}
	}
	photo := jpegWithOrientation(t, wide, 6)

	thumb, err := GenerateThumbnail(bytes.NewReader(photo), "image/jpeg")
	if err != nil {
		t.Fatal(err)
	}
	img, _, err := image.Decode(bytes.NewReader(thumb))
	if err != nil {
		t.Fatal(err)
	}
	if !isRed(img.At(100, 60)) || !isRed(img.At(300, 60)) || !isBlue(img.At(100, 340)) || !isBlue(img.At(300, 340)) {
		t.Errorf("the thumbnail is not the right way up: top %v %v, bottom %v %v",
			img.At(100, 60), img.At(300, 60), img.At(100, 340), img.At(300, 340))
	}

	processed, err := ProcessImage(bytes.NewReader(photo), "image/jpeg")
	if err != nil {
		t.Fatal(err)
	}
	if cfg, _, err := image.DecodeConfig(bytes.NewReader(processed)); err != nil || cfg.Width != 80 || cfg.Height != 160 {
		t.Errorf("ProcessImage = %dx%d (%v), want the photo turned upright at 80x160", cfg.Width, cfg.Height, err)
	}
}
`
}
