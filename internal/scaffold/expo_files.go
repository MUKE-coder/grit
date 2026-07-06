package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeExpoFiles(root string, opts Options) error {
	expoRoot := filepath.Join(root, "apps", "expo")

	files := map[string]string{
		filepath.Join(expoRoot, "package.json"):                            expoPackageJSON(opts),
		filepath.Join(expoRoot, "app.json"):                                expoAppJSON(opts),
		filepath.Join(expoRoot, "tsconfig.json"):                           expoTSConfig(),
		filepath.Join(expoRoot, "tailwind.config.js"):                      expoTailwindConfig(),
		filepath.Join(expoRoot, "metro.config.js"):                         expoMetroConfig(),
		filepath.Join(expoRoot, "babel.config.js"):                         expoBabelConfig(),
		filepath.Join(expoRoot, "global.css"):                              expoGlobalCSS(),
		filepath.Join(expoRoot, "nativewind-env.d.ts"):                     expoNativewindEnv(),
		filepath.Join(expoRoot, "app", "_layout.tsx"):                      expoRootLayout(),
		filepath.Join(expoRoot, "components", "ui", "pressable-scale.tsx"): expoPressableScale(),
		filepath.Join(expoRoot, "components", "ui", "screen-header.tsx"):   expoScreenHeader(),
		filepath.Join(expoRoot, "app", "(auth)", "_layout.tsx"):            expoAuthLayout(),
		filepath.Join(expoRoot, "app", "(auth)", "login.tsx"):              expoLoginScreen(),
		filepath.Join(expoRoot, "app", "(auth)", "register.tsx"):           expoRegisterScreen(),
		filepath.Join(expoRoot, "app", "(tabs)", "_layout.tsx"):            expoTabsLayout(),
		filepath.Join(expoRoot, "app", "(tabs)", "index.tsx"):              expoHomeScreen(),
		filepath.Join(expoRoot, "app", "(tabs)", "explore.tsx"):            expoExploreScreen(),
		filepath.Join(expoRoot, "app", "(tabs)", "profile.tsx"):            expoProfileScreen(),
		filepath.Join(expoRoot, "app", "(tabs)", "settings.tsx"):           expoSettingsScreen(),
		filepath.Join(expoRoot, "app", "explore", "users.tsx"):             expoExploreUsers(),
		filepath.Join(expoRoot, "app", "explore", "notifications.tsx"):     expoExploreNotifications(),
		filepath.Join(expoRoot, "app", "explore", "storage.tsx"):           expoExploreStorage(),
		filepath.Join(expoRoot, "app", "explore", "analytics.tsx"):         expoExploreAnalytics(),
		filepath.Join(expoRoot, "app", "explore", "content.tsx"):           expoExploreContent(),
		filepath.Join(expoRoot, "app", "explore", "integrations.tsx"):      expoExploreIntegrations(),
		filepath.Join(expoRoot, "app", "change-password.tsx"):              expoChangePasswordScreen(),
		filepath.Join(expoRoot, "app", "blogs", "index.tsx"):               expoBlogListScreen(),
		filepath.Join(expoRoot, "app", "blogs", "new.tsx"):                 expoBlogCreateScreen(),
		filepath.Join(expoRoot, "app", "users", "new.tsx"):                 expoUserCreateScreen(),
		filepath.Join(expoRoot, "hooks", "use-blogs.ts"):                   expoBlogHook(),
		filepath.Join(expoRoot, "lib", "upload.ts"):                        expoUploadHelper(),
		filepath.Join(expoRoot, "lib", "api.ts"):                           expoAPIClient(),
		filepath.Join(expoRoot, "lib", "auth.tsx"):                         expoAuthProvider(),
		filepath.Join(expoRoot, "lib", "theme.tsx"):                        expoThemeProvider(),
		filepath.Join(expoRoot, "lib", "query-client.ts"):                  expoQueryClient(),
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	// Icons/splash/favicon referenced by app.json — ship the Grit logo for each
	// so Metro doesn't fail with "Unable to resolve asset ./assets/icon.png".
	if err := writeBrandLogo(filepath.Join(expoRoot, "assets"),
		"icon.png", "splash.png", "adaptive-icon.png", "favicon.png"); err != nil {
		return err
	}

	return nil
}

func expoPackageJSON(opts Options) string {
	return fmt.Sprintf(`{
  "name": "%s-expo",
  "version": "1.0.0",
  "main": "expo-router/entry",
  "scripts": {
    "start": "expo start",
    "android": "expo start --android",
    "ios": "expo start --ios",
    "web": "expo start --web"
  },
  "dependencies": {
    "expo": "~54.0.0",
    "expo-router": "~6.0.0",
    "expo-linking": "~8.0.0",
    "expo-constants": "~18.0.0",
    "expo-status-bar": "~3.0.0",
    "expo-secure-store": "~15.0.0",
    "expo-splash-screen": "~31.0.0",
    "expo-image": "~3.0.11",
    "expo-haptics": "~15.0.8",
    "expo-web-browser": "~15.0.11",
    "expo-linear-gradient": "~15.0.7",
    "expo-blur": "~15.0.7",
    "expo-image-picker": "~17.0.8",
    "expo-file-system": "~19.0.15",
    "react": "19.1.0",
    "react-native": "0.81.5",
    "react-native-safe-area-context": "~5.6.0",
    "react-native-screens": "~4.16.0",
    "@react-navigation/bottom-tabs": "^7.0.0",
    "@expo/vector-icons": "^15.0.0",
    "nativewind": "^4.2.0",
    "react-native-css-interop": "^0.2.1",
    "@tanstack/react-query": "^5.0.0",
    "react-native-reanimated": "~4.1.0",
    "react-native-worklets": "0.5.1",
    "react-native-gesture-handler": "~2.28.0",
    "react-hook-form": "^7.54.0",
    "@hookform/resolvers": "^3.9.0",
    "zod": "^3.23.0"
  },
  "devDependencies": {
    "@babel/core": "^7.24.0",
    "@types/react": "~19.1.0",
    "tailwindcss": "^3.4.0",
    "typescript": "~5.9.2"
  },
  "private": true
}
`, opts.ProjectName)
}

func expoAppJSON(opts Options) string {
	return fmt.Sprintf(`{
  "expo": {
    "name": "%s",
    "slug": "%s",
    "version": "1.0.0",
    "orientation": "portrait",
    "scheme": "%s",
    "userInterfaceStyle": "dark",
    "newArchEnabled": true,
    "splash": {
      "image": "./assets/splash.png",
      "resizeMode": "contain",
      "backgroundColor": "#0a0a0f"
    },
    "icon": "./assets/icon.png",
    "ios": {
      "supportsTablet": true,
      "bundleIdentifier": "com.%s.app"
    },
    "android": {
      "adaptiveIcon": {
        "foregroundImage": "./assets/adaptive-icon.png",
        "backgroundColor": "#0a0a0f"
      },
      "package": "com.%s.app"
    },
    "web": {
      "bundler": "metro",
      "favicon": "./assets/favicon.png"
    },
    "plugins": [
      "expo-router",
      "expo-secure-store",
      [
        "expo-image-picker",
        {
          "photosPermission": "Allow this app to access your photos to set a profile picture."
        }
      ]
    ],
    "experiments": {
      "typedRoutes": true
    }
  }
}
`, opts.ProjectName, opts.ProjectName, opts.ProjectName, opts.ProjectName, opts.ProjectName)
}

func expoTSConfig() string {
	return `{
  "extends": "expo/tsconfig.base",
  "compilerOptions": {
    "strict": true,
    "paths": {
      "@/*": ["./*"]
    }
  },
  "include": ["**/*.ts", "**/*.tsx", ".expo/types/**/*.ts", "expo-env.d.ts", "nativewind-env.d.ts"]
}
`
}

func expoTailwindConfig() string {
	return `/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./app/**/*.{js,jsx,ts,tsx}", "./components/**/*.{js,jsx,ts,tsx}", "./lib/**/*.{js,jsx,ts,tsx}"],
  presets: [require("nativewind/preset")],
  // Class-based dark mode: the ThemeProvider drives it via NativeWind's
  // setColorScheme() so we control light/dark explicitly (default light),
  // instead of blindly following the OS setting.
  darkMode: "class",
  theme: {
    extend: {
      colors: {
        bg: {
          primary: "#0a0a0f",
          secondary: "#111118",
          tertiary: "#1a1a24",
          elevated: "#22222e",
          hover: "#2a2a38",
        },
        border: "#2a2a3a",
        accent: {
          DEFAULT: "#6c5ce7",
          hover: "#7c6cf7",
        },
        success: "#00b894",
        danger: "#ff6b6b",
        warning: "#fdcb6e",
        info: "#74b9ff",
      },
    },
  },
  plugins: [],
};
`
}

// expoBabelConfig is REQUIRED for NativeWind: without the jsxImportSource +
// nativewind/babel preset, every `className` is silently ignored and the app
// renders unstyled. The worklets plugin is required by react-native-reanimated
// v4 (its absence makes the app hang / crash on animated components).
func expoBabelConfig() string {
	return `module.exports = function (api) {
  api.cache(true);
  return {
    presets: [
      ["babel-preset-expo", { jsxImportSource: "nativewind" }],
      "nativewind/babel",
    ],
    plugins: ["react-native-worklets/plugin"],
  };
};
`
}

func expoMetroConfig() string {
	return `const { getDefaultConfig } = require("expo/metro-config");
const { withNativeWind } = require("nativewind/metro");

const config = getDefaultConfig(__dirname);

module.exports = withNativeWind(config, { input: "./global.css" });
`
}

func expoGlobalCSS() string {
	return `@tailwind base;
@tailwind components;
@tailwind utilities;
`
}

func expoNativewindEnv() string {
	return `/// <reference types="nativewind/types" />
`
}

func expoRootLayout() string {
	return `import "../global.css";
import { useEffect } from "react";
import { Stack, useRouter, useSegments } from "expo-router";
import { StatusBar } from "expo-status-bar";
import * as SplashScreen from "expo-splash-screen";
import { QueryClientProvider } from "@tanstack/react-query";
import { AuthProvider, useAuth } from "@/lib/auth";
import { ThemeProvider, useTheme } from "@/lib/theme";
import { queryClient } from "@/lib/query-client";

SplashScreen.preventAutoHideAsync();

function RootNav() {
  const { isAuthenticated, isLoading } = useAuth();
  const { palette } = useTheme();
  const router = useRouter();
  const segments = useSegments();

  // Declare BOTH groups and redirect imperatively. Conditionally rendering
  // one <Stack.Screen> sets the initial group but does NOT navigate when
  // auth flips after login — so a successful sign-in would leave you sitting
  // on the login screen. This effect moves you the moment auth changes.
  useEffect(() => {
    if (isLoading) return;
    const inAuthGroup = segments[0] === "(auth)";
    if (isAuthenticated && inAuthGroup) {
      router.replace("/(tabs)");
    } else if (!isAuthenticated && !inAuthGroup) {
      router.replace("/(auth)/login");
    }
  }, [isAuthenticated, isLoading, segments]);

  useEffect(() => {
    if (!isLoading) {
      SplashScreen.hideAsync();
    }
  }, [isLoading]);

  if (isLoading) return null;

  return (
    <>
      <StatusBar style={palette.statusBar} />
      <Stack screenOptions={{ headerShown: false }}>
        <Stack.Screen name="(auth)" />
        <Stack.Screen name="(tabs)" />
      </Stack>
    </>
  );
}

export default function RootLayout() {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider>
        <AuthProvider>
          <RootNav />
        </AuthProvider>
      </ThemeProvider>
    </QueryClientProvider>
  );
}
`
}

func expoPressableScale() string {
	return `import { forwardRef } from "react";
import { Pressable, type PressableProps, type ViewStyle } from "react-native";
import Animated, {
  useSharedValue,
  useAnimatedStyle,
  withSpring,
} from "react-native-reanimated";

const AnimatedPressable = Animated.createAnimatedComponent(Pressable);

interface PressableScaleProps extends PressableProps {
  /** Scale at full press. Defaults to 0.965 — subtle, iOS-app feel. */
  pressScale?: number;
  /** Tailwind class string (NativeWind). */
  className?: string;
}

/**
 * Press-state primitive with the spring-scale micro-interaction every
 * native app uses: scales down to ~96.5% on press-in, releases with a
 * critically-damped spring. Drop-in replacement for TouchableOpacity.
 */
export const PressableScale = forwardRef<typeof AnimatedPressable, PressableScaleProps>(
  function PressableScale(
    { pressScale = 0.965, style, onPressIn, onPressOut, children, ...rest },
    ref
  ) {
    const scale = useSharedValue(1);
    const animatedStyle = useAnimatedStyle(() => ({
      transform: [{ scale: scale.value }],
    }));

    return (
      <AnimatedPressable
        ref={ref as never}
        onPressIn={(e) => {
          scale.value = withSpring(pressScale, { damping: 18, stiffness: 320, mass: 0.6 });
          onPressIn?.(e);
        }}
        onPressOut={(e) => {
          scale.value = withSpring(1, { damping: 16, stiffness: 260, mass: 0.6 });
          onPressOut?.(e);
        }}
        style={[animatedStyle, style as ViewStyle]}
        {...rest}
      >
        {children as never}
      </AnimatedPressable>
    );
  }
);
`
}

// expoScreenHeader is the shared safe-area page header with an optional back
// button. Kept in sync with internal/generate mobileScreenHeaderContent so
// `grit generate resource` screens and hand-written pages look identical.
func expoScreenHeader() string {
	return `import { View, Text, Pressable } from "react-native";
import type { ReactNode } from "react";
import { useRouter } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { useTheme } from "@/lib/theme";

interface ScreenHeaderProps {
  title: string;
  subtitle?: string;
  showBack?: boolean;
  right?: ReactNode;
}

// Safe-area-aware page header. Non-tab screens pass showBack to get a back
// button; tab screens can use it without one for a consistent large title.
export function ScreenHeader({ title, subtitle, showBack = false, right }: ScreenHeaderProps) {
  const insets = useSafeAreaInsets();
  const router = useRouter();
  const { palette } = useTheme();
  return (
    <View style={{ paddingTop: insets.top + 8 }} className="px-6 pb-3 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <View className="flex-row items-center">
        {showBack ? (
          <Pressable
            onPress={() => router.back()}
            hitSlop={10}
            className="mr-3 w-9 h-9 rounded-full items-center justify-center bg-white dark:bg-[#1a1a24] border border-[#E5E7EB] dark:border-[#2a2a3a]"
          >
            <Ionicons name="chevron-back" size={20} color={palette.inputIcon} />
          </Pressable>
        ) : null}
        <View className="flex-1">
          <Text className="text-[26px] font-bold text-[#0F1018] dark:text-white tracking-tight">{title}</Text>
          {subtitle ? (
            <Text className="text-[14px] text-[#6B7280] dark:text-[#9090a8] mt-0.5">{subtitle}</Text>
          ) : null}
        </View>
        {right ? <View className="ml-3">{right}</View> : null}
      </View>
    </View>
  );
}
`
}

func expoAuthLayout() string {
	return `import { Stack } from "expo-router";
import { useTheme } from "@/lib/theme";

export default function AuthLayout() {
  const { scheme } = useTheme();
  return (
    <Stack
      screenOptions={{
        headerShown: false,
        contentStyle: { backgroundColor: scheme === "dark" ? "#0a0a0f" : "#F4F4F6" },
      }}
    />
  );
}
`
}

func expoLoginScreen() string {
	return `import { useState } from "react";
import {
  View,
  Text,
  TextInput,
  ActivityIndicator,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
  Pressable,
} from "react-native";
import { Link } from "expo-router";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Ionicons } from "@expo/vector-icons";
import { LinearGradient } from "expo-linear-gradient";
import { SafeAreaView } from "react-native-safe-area-context";
import Animated, { FadeInUp } from "react-native-reanimated";
import * as Haptics from "expo-haptics";
import { useAuth } from "@/lib/auth";
import { PressableScale } from "@/components/ui/pressable-scale";
import { Image } from "expo-image";
import { useTheme } from "@/lib/theme";

const loginSchema = z.object({
  email: z.string().email("Enter a valid email"),
  password: z.string().min(6, "Minimum 6 characters"),
});

type LoginForm = z.infer<typeof loginSchema>;

export default function LoginScreen() {
  const { login, loginWithGoogle } = useAuth();
  const [loading, setLoading] = useState(false);
  const [googleLoading, setGoogleLoading] = useState(false);
  const [error, setError] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const { palette } = useTheme();

  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: "", password: "" },
  });

  const onSubmit = async (data: LoginForm) => {
    setError("");
    setLoading(true);
    try {
      Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Light).catch(() => {});
      await login(data.email, data.password);
      Haptics.notificationAsync(Haptics.NotificationFeedbackType.Success).catch(() => {});
    } catch (err: any) {
      Haptics.notificationAsync(Haptics.NotificationFeedbackType.Error).catch(() => {});
      setError(err.message || "Login failed");
    } finally {
      setLoading(false);
    }
  };

  const handleGoogleLogin = async () => {
    setError("");
    setGoogleLoading(true);
    try {
      await loginWithGoogle();
    } catch (err: any) {
      setError(err.message || "Google login failed");
    } finally {
      setGoogleLoading(false);
    }
  };

  const emailBorder = errors.email ? "border-[#ff6b6b]" : "border-[#E5E7EB] dark:border-[#2a2a3a]";
  const passwordBorder = errors.password ? "border-[#ff6b6b]" : "border-[#E5E7EB] dark:border-[#2a2a3a]";

  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <FaintGrid />
      <SafeAreaView className="flex-1" edges={["top", "bottom"]}>
        <KeyboardAvoidingView
          behavior={Platform.OS === "ios" ? "padding" : "height"}
          keyboardVerticalOffset={Platform.OS === "ios" ? 0 : 24}
          className="flex-1"
        >
          <ScrollView
            contentContainerStyle={{ flexGrow: 1, justifyContent: "center", padding: 18 }}
            keyboardShouldPersistTaps="handled"
            automaticallyAdjustKeyboardInsets={Platform.OS === "ios"}
            showsVerticalScrollIndicator={false}
          >
            <Animated.View
              entering={FadeInUp.duration(500)}
              className="rounded-[28px] bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#1f1f2b] overflow-hidden"
              style={{
                shadowColor: "#6c5ce7",
                shadowOffset: { width: 0, height: 18 },
                shadowOpacity: 0.18,
                shadowRadius: 32,
                elevation: 8,
              }}
            >
              {/* Header — purple wash watermark */}
              <LinearGradient
                colors={palette.headerGradient}
                start={{ x: 0.5, y: 0 }}
                end={{ x: 0.5, y: 1 }}
                style={{ paddingTop: 36, paddingBottom: 16 }}
              >
                <View className="items-center">
                  <View
                    className="mb-5"
                    style={{
                      shadowColor: palette.logoShadow,
                      shadowOffset: { width: 0, height: 14 },
                      shadowOpacity: 0.35,
                      shadowRadius: 22,
                    }}
                  >
                    <Image
                      source={require("../../assets/icon.png")}
                      style={{ width: 76, height: 76, borderRadius: 20 }}
                      contentFit="contain"
                    />
                  </View>
                  <Text className="text-[28px] font-bold text-[#0F1018] dark:text-white tracking-tight">
                    Welcome back
                  </Text>
                  <Text className="text-[#6B7280] dark:text-[#9090a8] text-[14px] mt-1.5 px-6 text-center leading-5">
                    Sign in to your account to continue.
                  </Text>
                </View>
              </LinearGradient>

              {/* Form body */}
              <View className="px-7 pt-7 pb-7">
                {error ? (
                  <Animated.View
                    entering={FadeInUp.duration(280)}
                    className="bg-[#ff6b6b]/10 border border-[#ff6b6b]/25 rounded-2xl px-4 py-3 mb-4 flex-row items-center"
                  >
                    <Ionicons name="alert-circle" size={18} color="#ff6b6b" />
                    <Text className="text-[#ff6b6b] text-[13px] ml-2 flex-1">{error}</Text>
                  </Animated.View>
                ) : null}

                {/* Email */}
                <View className="mb-4">
                  <Text className="text-[12.5px] font-semibold text-[#6B7280] dark:text-[#9090a8] mb-2 tracking-tight">
                    Email Address
                  </Text>
                  <Controller
                    control={control}
                    name="email"
                    render={({ field: { onChange, onBlur, value } }) => (
                      <View
                        className={"bg-[#F4F4F6] dark:bg-[#0a0a0f] border rounded-2xl flex-row items-center px-4 " + emailBorder}
                        style={{ height: 52 }}
                      >
                        <Ionicons name="mail-outline" size={17} color="#606078" />
                        <TextInput
                          className="flex-1 ml-2.5 text-[#0F1018] dark:text-white text-[15px]"
                          placeholder="you@example.com"
                          placeholderTextColor="#606078"
                          value={value}
                          onChangeText={onChange}
                          onBlur={onBlur}
                          keyboardType="email-address"
                          autoCapitalize="none"
                          autoComplete="email"
                          autoCorrect={false}
                        />
                      </View>
                    )}
                  />
                  {errors.email ? (
                    <Text className="text-[#ff6b6b] text-[11.5px] mt-1.5 ml-1">
                      {errors.email.message}
                    </Text>
                  ) : null}
                </View>

                {/* Password */}
                <View className="mb-6">
                  <Text className="text-[12.5px] font-semibold text-[#6B7280] dark:text-[#9090a8] mb-2 tracking-tight">
                    Password
                  </Text>
                  <Controller
                    control={control}
                    name="password"
                    render={({ field: { onChange, onBlur, value } }) => (
                      <View
                        className={"bg-[#F4F4F6] dark:bg-[#0a0a0f] border rounded-2xl flex-row items-center px-4 " + passwordBorder}
                        style={{ height: 52 }}
                      >
                        <Ionicons name="lock-closed-outline" size={17} color="#606078" />
                        <TextInput
                          className="flex-1 ml-2.5 text-[#0F1018] dark:text-white text-[15px]"
                          placeholder="••••••••"
                          placeholderTextColor="#606078"
                          value={value}
                          onChangeText={onChange}
                          onBlur={onBlur}
                          secureTextEntry={!showPassword}
                          autoCapitalize="none"
                          autoCorrect={false}
                        />
                        <Pressable
                          onPress={() => {
                            Haptics.selectionAsync().catch(() => {});
                            setShowPassword((s) => !s);
                          }}
                          hitSlop={10}
                          className="ml-2 p-1"
                        >
                          <Ionicons
                            name={showPassword ? "eye-off-outline" : "eye-outline"}
                            size={19}
                            color={showPassword ? "#6c5ce7" : "#606078"}
                          />
                        </Pressable>
                      </View>
                    )}
                  />
                  {errors.password ? (
                    <Text className="text-[#ff6b6b] text-[11.5px] mt-1.5 ml-1">
                      {errors.password.message}
                    </Text>
                  ) : null}
                </View>

                {/* Primary CTA */}
                <PressableScale
                  onPress={handleSubmit(onSubmit)}
                  disabled={loading || googleLoading}
                  pressScale={0.97}
                  className="overflow-hidden rounded-full"
                  style={{
                    shadowColor: "#6c5ce7",
                    shadowOffset: { width: 0, height: 10 },
                    shadowOpacity: 0.4,
                    shadowRadius: 18,
                    opacity: loading ? 0.85 : 1,
                  }}
                >
                  <LinearGradient
                    colors={["#7c6cf7", "#6c5ce7"]}
                    start={{ x: 0, y: 0 }}
                    end={{ x: 1, y: 1 }}
                    style={{ height: 54, flexDirection: "row", alignItems: "center", justifyContent: "center" }}
                  >
                    {loading ? (
                      <ActivityIndicator color="#fff" />
                    ) : (
                      <>
                        <Text className="text-white text-[15.5px] font-semibold tracking-tight">
                          Sign in
                        </Text>
                        <Ionicons name="arrow-forward" size={17} color="#FFFFFF" style={{ marginLeft: 8 }} />
                      </>
                    )}
                  </LinearGradient>
                </PressableScale>

                {/* Divider */}
                <View className="flex-row items-center my-6">
                  <View className="flex-1 h-px bg-[#E5E7EB] dark:bg-[#2a2a3a]" />
                  <Text className="text-[#9CA3AF] dark:text-[#606078] mx-4 text-[12px]">or</Text>
                  <View className="flex-1 h-px bg-[#E5E7EB] dark:bg-[#2a2a3a]" />
                </View>

                {/* Google */}
                <PressableScale
                  onPress={handleGoogleLogin}
                  disabled={loading || googleLoading}
                  className="rounded-full border border-[#E5E7EB] dark:border-[#2a2a3a] bg-[#F4F4F6] dark:bg-[#0a0a0f] flex-row items-center justify-center"
                  style={{ height: 52 }}
                >
                  {googleLoading ? (
                    <ActivityIndicator color="#fff" />
                  ) : (
                    <>
                      <Ionicons name="logo-google" size={19} color="#e8e8f0" />
                      <Text className="text-[#0F1018] dark:text-white font-semibold text-[15px] ml-3">
                        Continue with Google
                      </Text>
                    </>
                  )}
                </PressableScale>

                <View className="flex-row justify-center mt-6">
                  <Text className="text-[#6B7280] dark:text-[#9090a8] text-[13.5px]">Don't have an account? </Text>
                  <Link href="/(auth)/register">
                    <Text className="text-[#6c5ce7] font-semibold text-[13.5px]">Sign up</Text>
                  </Link>
                </View>
              </View>
            </Animated.View>
          </ScrollView>
        </KeyboardAvoidingView>
      </SafeAreaView>
    </View>
  );
}

// Subtle architectural grid behind the auth flow — ties login and
// register into one visual story.
function FaintGrid() {
  const { palette } = useTheme();
  return (
    <View
      style={{ position: "absolute", top: 0, left: 0, right: 0, bottom: 0, opacity: palette.gridOpacity }}
      pointerEvents="none"
    >
      <View style={{ flexDirection: "row", justifyContent: "space-between", position: "absolute", top: 0, left: 0, right: 0, bottom: 0 }}>
        {Array.from({ length: 10 }).map((_, i) => (
          <View key={"gv-" + i} style={{ width: 1, backgroundColor: palette.gridLine }} />
        ))}
      </View>
      <View style={{ flexDirection: "column", justifyContent: "space-between", position: "absolute", top: 0, left: 0, right: 0, bottom: 0 }}>
        {Array.from({ length: 18 }).map((_, i) => (
          <View key={"gh-" + i} style={{ height: 1, backgroundColor: palette.gridLine }} />
        ))}
      </View>
    </View>
  );
}
`
}

func expoRegisterScreen() string {
	return `import { useState } from "react";
import {
  View,
  Text,
  TextInput,
  ActivityIndicator,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
  Pressable,
} from "react-native";
import { Link } from "expo-router";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Ionicons } from "@expo/vector-icons";
import { LinearGradient } from "expo-linear-gradient";
import { SafeAreaView } from "react-native-safe-area-context";
import Animated, { FadeInUp } from "react-native-reanimated";
import * as Haptics from "expo-haptics";
import { useAuth } from "@/lib/auth";
import { PressableScale } from "@/components/ui/pressable-scale";
import { Image } from "expo-image";
import { useTheme } from "@/lib/theme";

const registerSchema = z.object({
  firstName: z.string().min(1, "Required"),
  lastName: z.string().min(1, "Required"),
  email: z.string().email("Enter a valid email"),
  password: z.string().min(8, "Minimum 8 characters"),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords do not match",
  path: ["confirmPassword"],
});

type RegisterForm = z.infer<typeof registerSchema>;

export default function RegisterScreen() {
  const { register: registerUser } = useAuth();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const { palette } = useTheme();

  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterForm>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      firstName: "",
      lastName: "",
      email: "",
      password: "",
      confirmPassword: "",
    },
  });

  const onSubmit = async (data: RegisterForm) => {
    setError("");
    setLoading(true);
    try {
      Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Light).catch(() => {});
      await registerUser(data.firstName, data.lastName, data.email, data.password);
      Haptics.notificationAsync(Haptics.NotificationFeedbackType.Success).catch(() => {});
    } catch (err: any) {
      Haptics.notificationAsync(Haptics.NotificationFeedbackType.Error).catch(() => {});
      setError(err.message || "Registration failed");
    } finally {
      setLoading(false);
    }
  };

  const inputBase = "bg-[#F4F4F6] dark:bg-[#0a0a0f] border rounded-2xl flex-row items-center px-4";
  const border = (hasError?: boolean) => (hasError ? "border-[#ff6b6b]" : "border-[#E5E7EB] dark:border-[#2a2a3a]");
  const label = "text-[12.5px] font-semibold text-[#6B7280] dark:text-[#9090a8] mb-2 tracking-tight";
  const errorText = "text-[#ff6b6b] text-[11.5px] mt-1.5 ml-1";

  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <FaintGrid />
      <SafeAreaView className="flex-1" edges={["top", "bottom"]}>
        <KeyboardAvoidingView
          behavior={Platform.OS === "ios" ? "padding" : "height"}
          keyboardVerticalOffset={Platform.OS === "ios" ? 0 : 24}
          className="flex-1"
        >
          <ScrollView
            contentContainerStyle={{ flexGrow: 1, justifyContent: "center", padding: 18 }}
            keyboardShouldPersistTaps="handled"
            automaticallyAdjustKeyboardInsets={Platform.OS === "ios"}
            showsVerticalScrollIndicator={false}
          >
            <Animated.View
              entering={FadeInUp.duration(500)}
              className="rounded-[28px] bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#1f1f2b] overflow-hidden"
              style={{
                shadowColor: "#6c5ce7",
                shadowOffset: { width: 0, height: 18 },
                shadowOpacity: 0.18,
                shadowRadius: 32,
                elevation: 8,
              }}
            >
              <LinearGradient
                colors={palette.headerGradient}
                start={{ x: 0.5, y: 0 }}
                end={{ x: 0.5, y: 1 }}
                style={{ paddingTop: 36, paddingBottom: 16 }}
              >
                <View className="items-center">
                  <View
                    className="mb-5"
                    style={{
                      shadowColor: palette.logoShadow,
                      shadowOffset: { width: 0, height: 14 },
                      shadowOpacity: 0.35,
                      shadowRadius: 22,
                    }}
                  >
                    <Image
                      source={require("../../assets/icon.png")}
                      style={{ width: 76, height: 76, borderRadius: 20 }}
                      contentFit="contain"
                    />
                  </View>
                  <Text className="text-[28px] font-bold text-[#0F1018] dark:text-white tracking-tight">
                    Create account
                  </Text>
                  <Text className="text-[#6B7280] dark:text-[#9090a8] text-[14px] mt-1.5 px-6 text-center leading-5">
                    Get started with Grit in seconds.
                  </Text>
                </View>
              </LinearGradient>

              <View className="px-7 pt-7 pb-7">
                {error ? (
                  <Animated.View
                    entering={FadeInUp.duration(280)}
                    className="bg-[#ff6b6b]/10 border border-[#ff6b6b]/25 rounded-2xl px-4 py-3 mb-4 flex-row items-center"
                  >
                    <Ionicons name="alert-circle" size={18} color="#ff6b6b" />
                    <Text className="text-[#ff6b6b] text-[13px] ml-2 flex-1">{error}</Text>
                  </Animated.View>
                ) : null}

                <View className="flex-row gap-3 mb-4">
                  <View className="flex-1">
                    <Text className={label}>First name</Text>
                    <Controller
                      control={control}
                      name="firstName"
                      render={({ field: { onChange, onBlur, value } }) => (
                        <View className={inputBase + " " + border(!!errors.firstName)} style={{ height: 52 }}>
                          <TextInput
                            className="flex-1 text-[#0F1018] dark:text-white text-[15px]"
                            placeholder="John"
                            placeholderTextColor="#606078"
                            value={value}
                            onChangeText={onChange}
                            onBlur={onBlur}
                            autoComplete="given-name"
                          />
                        </View>
                      )}
                    />
                    {errors.firstName ? (
                      <Text className={errorText}>{errors.firstName.message}</Text>
                    ) : null}
                  </View>
                  <View className="flex-1">
                    <Text className={label}>Last name</Text>
                    <Controller
                      control={control}
                      name="lastName"
                      render={({ field: { onChange, onBlur, value } }) => (
                        <View className={inputBase + " " + border(!!errors.lastName)} style={{ height: 52 }}>
                          <TextInput
                            className="flex-1 text-[#0F1018] dark:text-white text-[15px]"
                            placeholder="Doe"
                            placeholderTextColor="#606078"
                            value={value}
                            onChangeText={onChange}
                            onBlur={onBlur}
                            autoComplete="family-name"
                          />
                        </View>
                      )}
                    />
                    {errors.lastName ? (
                      <Text className={errorText}>{errors.lastName.message}</Text>
                    ) : null}
                  </View>
                </View>

                <View className="mb-4">
                  <Text className={label}>Email Address</Text>
                  <Controller
                    control={control}
                    name="email"
                    render={({ field: { onChange, onBlur, value } }) => (
                      <View className={inputBase + " " + border(!!errors.email)} style={{ height: 52 }}>
                        <Ionicons name="mail-outline" size={17} color="#606078" />
                        <TextInput
                          className="flex-1 ml-2.5 text-[#0F1018] dark:text-white text-[15px]"
                          placeholder="you@example.com"
                          placeholderTextColor="#606078"
                          value={value}
                          onChangeText={onChange}
                          onBlur={onBlur}
                          keyboardType="email-address"
                          autoCapitalize="none"
                          autoComplete="email"
                          autoCorrect={false}
                        />
                      </View>
                    )}
                  />
                  {errors.email ? <Text className={errorText}>{errors.email.message}</Text> : null}
                </View>

                <View className="mb-4">
                  <Text className={label}>Password</Text>
                  <Controller
                    control={control}
                    name="password"
                    render={({ field: { onChange, onBlur, value } }) => (
                      <View className={inputBase + " " + border(!!errors.password)} style={{ height: 52 }}>
                        <Ionicons name="lock-closed-outline" size={17} color="#606078" />
                        <TextInput
                          className="flex-1 ml-2.5 text-[#0F1018] dark:text-white text-[15px]"
                          placeholder="Min. 8 characters"
                          placeholderTextColor="#606078"
                          value={value}
                          onChangeText={onChange}
                          onBlur={onBlur}
                          secureTextEntry={!showPassword}
                          autoCapitalize="none"
                          autoCorrect={false}
                        />
                        <Pressable
                          onPress={() => {
                            Haptics.selectionAsync().catch(() => {});
                            setShowPassword((s) => !s);
                          }}
                          hitSlop={10}
                          className="ml-2 p-1"
                        >
                          <Ionicons
                            name={showPassword ? "eye-off-outline" : "eye-outline"}
                            size={19}
                            color={showPassword ? "#6c5ce7" : "#606078"}
                          />
                        </Pressable>
                      </View>
                    )}
                  />
                  {errors.password ? <Text className={errorText}>{errors.password.message}</Text> : null}
                </View>

                <View className="mb-6">
                  <Text className={label}>Confirm password</Text>
                  <Controller
                    control={control}
                    name="confirmPassword"
                    render={({ field: { onChange, onBlur, value } }) => (
                      <View className={inputBase + " " + border(!!errors.confirmPassword)} style={{ height: 52 }}>
                        <Ionicons name="lock-closed-outline" size={17} color="#606078" />
                        <TextInput
                          className="flex-1 ml-2.5 text-[#0F1018] dark:text-white text-[15px]"
                          placeholder="Repeat password"
                          placeholderTextColor="#606078"
                          value={value}
                          onChangeText={onChange}
                          onBlur={onBlur}
                          secureTextEntry={!showPassword}
                          autoCapitalize="none"
                          autoCorrect={false}
                        />
                      </View>
                    )}
                  />
                  {errors.confirmPassword ? (
                    <Text className={errorText}>{errors.confirmPassword.message}</Text>
                  ) : null}
                </View>

                <PressableScale
                  onPress={handleSubmit(onSubmit)}
                  disabled={loading}
                  pressScale={0.97}
                  className="overflow-hidden rounded-full"
                  style={{
                    shadowColor: "#6c5ce7",
                    shadowOffset: { width: 0, height: 10 },
                    shadowOpacity: 0.4,
                    shadowRadius: 18,
                    opacity: loading ? 0.85 : 1,
                  }}
                >
                  <LinearGradient
                    colors={["#7c6cf7", "#6c5ce7"]}
                    start={{ x: 0, y: 0 }}
                    end={{ x: 1, y: 1 }}
                    style={{ height: 54, flexDirection: "row", alignItems: "center", justifyContent: "center" }}
                  >
                    {loading ? (
                      <ActivityIndicator color="#fff" />
                    ) : (
                      <>
                        <Text className="text-white text-[15.5px] font-semibold tracking-tight">
                          Create account
                        </Text>
                        <Ionicons name="arrow-forward" size={17} color="#FFFFFF" style={{ marginLeft: 8 }} />
                      </>
                    )}
                  </LinearGradient>
                </PressableScale>

                <View className="flex-row justify-center mt-6">
                  <Text className="text-[#6B7280] dark:text-[#9090a8] text-[13.5px]">Already have an account? </Text>
                  <Link href="/(auth)/login">
                    <Text className="text-[#6c5ce7] font-semibold text-[13.5px]">Sign in</Text>
                  </Link>
                </View>
              </View>
            </Animated.View>
          </ScrollView>
        </KeyboardAvoidingView>
      </SafeAreaView>
    </View>
  );
}

function FaintGrid() {
  const { palette } = useTheme();
  return (
    <View
      style={{ position: "absolute", top: 0, left: 0, right: 0, bottom: 0, opacity: palette.gridOpacity }}
      pointerEvents="none"
    >
      <View style={{ flexDirection: "row", justifyContent: "space-between", position: "absolute", top: 0, left: 0, right: 0, bottom: 0 }}>
        {Array.from({ length: 10 }).map((_, i) => (
          <View key={"gv-" + i} style={{ width: 1, backgroundColor: palette.gridLine }} />
        ))}
      </View>
      <View style={{ flexDirection: "column", justifyContent: "space-between", position: "absolute", top: 0, left: 0, right: 0, bottom: 0 }}>
        {Array.from({ length: 18 }).map((_, i) => (
          <View key={"gh-" + i} style={{ height: 1, backgroundColor: palette.gridLine }} />
        ))}
      </View>
    </View>
  );
}
`
}

func expoTabsLayout() string {
	return `import { Tabs } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import { BlurView } from "expo-blur";
import { Platform, StyleSheet, View } from "react-native";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import * as Haptics from "expo-haptics";
import { useTheme } from "@/lib/theme";

// Floating glass tab bar. iOS gets a native frosted-blur background;
// Android falls back to a solid surface. Colours follow the active theme.
export default function TabsLayout() {
  const { palette, scheme } = useTheme();
  const insets = useSafeAreaInsets();
  // Lift the floating bar above the OS gesture/nav area so it isn't clipped
  // by the Android system navigation at the very bottom of the screen.
  const barBottom = Math.max(insets.bottom, 10) + 6;
  const blurBg =
    scheme === "dark" ? "rgba(17,17,24,0.6)" : "rgba(255,255,255,0.7)";
  return (
    <Tabs
      screenListeners={{
        tabPress: () => {
          Haptics.selectionAsync().catch(() => {});
        },
      }}
      screenOptions={{
        headerShown: false,
        tabBarActiveTintColor: "#7c6cf7",
        tabBarInactiveTintColor: palette.tabInactive,
        tabBarLabelStyle: {
          fontSize: 11,
          fontWeight: "600",
        },
        tabBarItemStyle: {
          paddingTop: 6,
        },
        tabBarBackground:
          Platform.OS === "ios"
            ? () => (
                <BlurView
                  tint={palette.blurTint}
                  intensity={40}
                  style={[
                    StyleSheet.absoluteFill,
                    { backgroundColor: blurBg, borderRadius: 24, overflow: "hidden" },
                  ]}
                />
              )
            : () => (
                <View
                  style={[
                    StyleSheet.absoluteFill,
                    { backgroundColor: palette.tabBar, borderRadius: 24 },
                  ]}
                />
              ),
        tabBarStyle: {
          position: "absolute",
          left: 16,
          right: 16,
          bottom: barBottom,
          height: 64,
          paddingBottom: 8,
          borderRadius: 24,
          borderTopWidth: 0,
          borderWidth: 1,
          borderColor: palette.tabBarBorder,
          backgroundColor: "transparent",
          elevation: 12,
          shadowColor: "#000",
          shadowOffset: { width: 0, height: 8 },
          shadowOpacity: scheme === "dark" ? 0.35 : 0.12,
          shadowRadius: 16,
        },
      }}
    >
      <Tabs.Screen
        name="index"
        options={{
          title: "Home",
          tabBarIcon: ({ color, size }) => (
            <Ionicons name="grid-outline" size={size} color={color} />
          ),
        }}
      />
      <Tabs.Screen
        name="explore"
        options={{
          title: "More",
          tabBarIcon: ({ color, size }) => (
            <Ionicons name="ellipsis-horizontal" size={size} color={color} />
          ),
        }}
      />
      <Tabs.Screen
        name="profile"
        options={{
          title: "Profile",
          tabBarIcon: ({ color, size }) => (
            <Ionicons name="person-outline" size={size} color={color} />
          ),
        }}
      />
      <Tabs.Screen
        name="settings"
        options={{
          title: "Settings",
          tabBarIcon: ({ color, size }) => (
            <Ionicons name="settings-outline" size={size} color={color} />
          ),
        }}
      />
    </Tabs>
  );
}
`
}

func expoHomeScreen() string {
	return `import { View, Text, ScrollView, RefreshControl, FlatList } from "react-native";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import { useState, useCallback } from "react";
import { Ionicons } from "@expo/vector-icons";

interface Stats {
  total_users: number;
  active_users: number;
  new_today: number;
  total_items: number;
}

interface RecentItem {
  id: string;
  title: string;
  subtitle: string;
  icon: string;
  time: string;
}

function StatCard({
  title,
  value,
  color,
  icon,
}: {
  title: string;
  value: number;
  color: string;
  icon: string;
}) {
  return (
    <View className="bg-white dark:bg-[#22222e] border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-2xl p-4 flex-1 min-w-[140px]">
      <View className="flex-row items-center justify-between mb-2">
        <Ionicons name={icon as any} size={18} color={color} />
      </View>
      <Text className="text-2xl font-bold text-[#0F1018] dark:text-white">{value}</Text>
      <Text className="text-xs text-[#6B7280] dark:text-[#9090a8] mt-1">{title}</Text>
    </View>
  );
}

function RecentItemRow({ item }: { item: RecentItem }) {
  return (
    <View className="flex-row items-center bg-white dark:bg-[#22222e] border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-xl px-4 py-3 mb-2">
      <View className="w-10 h-10 rounded-full bg-[#6c5ce7]/20 items-center justify-center mr-3">
        <Ionicons name={item.icon as any} size={18} color="#6c5ce7" />
      </View>
      <View className="flex-1">
        <Text className="text-sm font-medium text-[#0F1018] dark:text-white">{item.title}</Text>
        <Text className="text-xs text-[#6B7280] dark:text-[#9090a8] mt-0.5">{item.subtitle}</Text>
      </View>
      <Text className="text-xs text-[#9CA3AF] dark:text-[#606078]">{item.time}</Text>
    </View>
  );
}

export default function HomeScreen() {
  const { user } = useAuth();
  const [refreshing, setRefreshing] = useState(false);

  const { data: stats, refetch } = useQuery<Stats>({
    queryKey: ["home-stats"],
    queryFn: async () => {
      const res = await api.get("/users?page=1&page_size=1");
      return {
        total_users: res.data?.meta?.total || 0,
        active_users: res.data?.meta?.total || 0,
        new_today: 0,
        total_items: 0,
      };
    },
  });

  const { data: recentItems } = useQuery<RecentItem[]>({
    queryKey: ["recent-items"],
    queryFn: async () => {
      // Default placeholder items until the API provides recent activity
      return [
        { id: "1", title: "App launched", subtitle: "Your project is running", icon: "rocket-outline", time: "Now" },
        { id: "2", title: "API connected", subtitle: "Backend is reachable", icon: "cloud-done-outline", time: "Now" },
        { id: "3", title: "Auth ready", subtitle: "Login and register work", icon: "shield-checkmark-outline", time: "Now" },
      ];
    },
  });

  const onRefresh = useCallback(async () => {
    setRefreshing(true);
    await refetch();
    setRefreshing(false);
  }, [refetch]);

  const firstName = user?.first_name || user?.name?.split(" ")[0] || "User";

  return (
    <ScrollView
      className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]"
      contentContainerClassName="px-6 pt-16 pb-28"
      refreshControl={
        <RefreshControl refreshing={refreshing} onRefresh={onRefresh} tintColor="#6c5ce7" />
      }
    >
      <Text className="text-2xl font-bold text-[#0F1018] dark:text-white mb-1">
        Welcome back, {firstName}
      </Text>
      <Text className="text-base text-[#6B7280] dark:text-[#9090a8] mb-6">
        Here's what's happening today.
      </Text>

      <View className="flex-row gap-3 mb-6">
        <StatCard title="Total Users" value={stats?.total_users || 0} color="#6c5ce7" icon="people-outline" />
        <StatCard title="Active" value={stats?.active_users || 0} color="#00b894" icon="pulse-outline" />
      </View>

      <View className="flex-row gap-3 mb-8">
        <StatCard title="New Today" value={stats?.new_today || 0} color="#74b9ff" icon="trending-up-outline" />
        <StatCard title="Items" value={stats?.total_items || 0} color="#fdcb6e" icon="cube-outline" />
      </View>

      <Text className="text-lg font-semibold text-[#0F1018] dark:text-white mb-3">Recent Activity</Text>

      {recentItems?.map((item) => (
        <RecentItemRow key={item.id} item={item} />
      ))}

      <View className="bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-2xl p-6 mt-4">
        <Text className="text-lg font-semibold text-[#0F1018] dark:text-white mb-3">Quick Start</Text>
        <Text className="text-sm text-[#6B7280] dark:text-[#9090a8] leading-6">
          Your Grit mobile app is connected to the API. Edit this screen in{"\n"}
          apps/expo/app/(tabs)/index.tsx
        </Text>
      </View>
    </ScrollView>
  );
}
`
}

func expoExploreScreen() string {
	return `import { View, Text, ScrollView, TouchableOpacity } from "react-native";
import { useRouter } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import { ScreenHeader } from "@/components/ui/screen-header";

interface LinkItem {
  title: string;
  description: string;
  icon: string;
  color: string;
  route: string;
}

// Generated resources. ` + "`" + `grit generate resource` + "`" + ` injects an entry below the
// marker for every resource, so each one's list screen is reachable here.
const resources: LinkItem[] = [
  { title: "Users", description: "Manage user accounts", icon: "people-outline", color: "#6c5ce7", route: "/explore/users" },
  { title: "Blogs", description: "Posts and articles", icon: "newspaper-outline", color: "#00b894", route: "/blogs" },
  // grit:mobile-resources
];

const tools: LinkItem[] = [
  { title: "Content", description: "Posts, pages, and media", icon: "document-text-outline", color: "#00b894", route: "/explore/content" },
  { title: "Analytics", description: "Usage and performance", icon: "bar-chart-outline", color: "#74b9ff", route: "/explore/analytics" },
  { title: "Notifications", description: "Alerts and messages", icon: "notifications-outline", color: "#fdcb6e", route: "/explore/notifications" },
  { title: "Storage", description: "Files and uploads", icon: "cloud-outline", color: "#ff6b6b", route: "/explore/storage" },
  { title: "Integrations", description: "Connected services", icon: "extension-puzzle-outline", color: "#a29bfe", route: "/explore/integrations" },
];

function LinkCard({ item }: { item: LinkItem }) {
  const router = useRouter();
  return (
    <TouchableOpacity
      className="bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#1f1f2b] rounded-2xl p-5 mb-3"
      activeOpacity={0.7}
      onPress={() => router.push(item.route as any)}
    >
      <View className="flex-row items-center">
        <View
          className="w-11 h-11 rounded-xl items-center justify-center mr-4"
          style={{ backgroundColor: item.color + "20" }}
        >
          <Ionicons name={item.icon as any} size={22} color={item.color} />
        </View>
        <View className="flex-1">
          <Text className="text-base font-semibold text-[#0F1018] dark:text-white">{item.title}</Text>
          <Text className="text-xs text-[#6B7280] dark:text-[#9090a8] mt-0.5">{item.description}</Text>
        </View>
        <Ionicons name="chevron-forward" size={18} color="#9CA3AF" />
      </View>
    </TouchableOpacity>
  );
}

function SectionTitle({ children }: { children: string }) {
  return (
    <Text className="text-[13px] font-semibold uppercase tracking-wider text-[#9CA3AF] dark:text-[#606078] mb-3 mt-2">
      {children}
    </Text>
  );
}

export default function MoreScreen() {
  return (
    <View className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
      <ScreenHeader title="More" subtitle="Resources & tools" />
      <ScrollView contentContainerStyle={{ paddingHorizontal: 24, paddingBottom: 120 }}>
        <SectionTitle>Resources</SectionTitle>
        {resources.length > 0 ? (
          resources.map((item) => <LinkCard key={item.route} item={item} />)
        ) : (
          <View className="bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#1f1f2b] rounded-2xl p-5 mb-3">
            <Text className="text-[14px] text-[#6B7280] dark:text-[#9090a8] leading-5">
              Generate a resource and it shows up here:{"\n"}
              <Text className="font-semibold text-[#6c5ce7]">grit generate resource Product</Text>
            </Text>
          </View>
        )}

        <SectionTitle>Tools</SectionTitle>
        {tools.map((item) => (
          <LinkCard key={item.route} item={item} />
        ))}
      </ScrollView>
    </View>
  );
}
`
}

func expoProfileScreen() string {
	return `import { useState } from "react";
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  ScrollView,
  ActivityIndicator,
  Alert,
} from "react-native";
import { Image } from "expo-image";
import { useRouter } from "expo-router";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useAuth } from "@/lib/auth";
import { api } from "@/lib/api";
import { pickAndUploadImage } from "@/lib/upload";
import { Ionicons } from "@expo/vector-icons";

const profileSchema = z.object({
  firstName: z.string().min(1, "Required"),
  lastName: z.string().min(1, "Required"),
  email: z.string().email("Invalid email"),
});

type ProfileForm = z.infer<typeof profileSchema>;

function ProfileRow({
  icon,
  label,
  value,
}: {
  icon: string;
  label: string;
  value: string;
}) {
  return (
    <View className="flex-row items-center px-5 py-4">
      <Ionicons name={icon as any} size={20} color="#9090a8" />
      <Text className="text-sm text-[#6B7280] dark:text-[#9090a8] ml-3 flex-1">{label}</Text>
      <Text className="text-sm text-[#0F1018] dark:text-white">{value}</Text>
    </View>
  );
}

export default function ProfileScreen() {
  const { user, logout, refreshUser } = useAuth();
  const router = useRouter();
  const [editing, setEditing] = useState(false);
  const [saving, setSaving] = useState(false);
  const [uploadingAvatar, setUploadingAvatar] = useState(false);
  const [error, setError] = useState("");

  const onChangeAvatar = async () => {
    try {
      setUploadingAvatar(true);
      const url = await pickAndUploadImage();
      if (url) {
        await api.put("/profile", { avatar: url });
        await refreshUser();
      }
    } catch (err: any) {
      Alert.alert("Upload failed", err.message || "Could not update your photo");
    } finally {
      setUploadingAvatar(false);
    }
  };

  const displayName =
    (user?.first_name || "") + " " + (user?.last_name || "") || user?.name || "User";
  const initials =
    (user?.first_name?.charAt(0) || user?.name?.charAt(0) || "U").toUpperCase();

  const {
    control,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<ProfileForm>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      firstName: user?.first_name || "",
      lastName: user?.last_name || "",
      email: user?.email || "",
    },
  });

  const onSave = async (data: ProfileForm) => {
    setError("");
    setSaving(true);
    try {
      await api.put("/profile", {
        first_name: data.firstName,
        last_name: data.lastName,
        email: data.email,
      });
      await refreshUser();
      setEditing(false);
    } catch (err: any) {
      setError(err.message || "Update failed");
    } finally {
      setSaving(false);
    }
  };

  const onCancel = () => {
    reset();
    setError("");
    setEditing(false);
  };

  if (editing) {
    return (
      <ScrollView className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]" contentContainerClassName="px-6 pt-16 pb-28">
        <Text className="text-2xl font-bold text-[#0F1018] dark:text-white mb-6">Edit Profile</Text>

        {error ? (
          <View className="bg-[#ff6b6b]/10 border border-[#ff6b6b]/30 rounded-xl p-4 mb-6">
            <Text className="text-[#ff6b6b] text-sm">{error}</Text>
          </View>
        ) : null}

        <View className="mb-4">
          <Text className="text-sm text-[#6B7280] dark:text-[#9090a8] mb-2">First name</Text>
          <Controller
            control={control}
            name="firstName"
            render={({ field: { onChange, onBlur, value } }) => (
              <TextInput
                className="bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-xl px-4 py-3.5 text-[#0F1018] dark:text-white text-base"
                placeholder="First name"
                placeholderTextColor="#606078"
                value={value}
                onChangeText={onChange}
                onBlur={onBlur}
              />
            )}
          />
          {errors.firstName ? (
            <Text className="text-[#ff6b6b] text-xs mt-1">{errors.firstName.message}</Text>
          ) : null}
        </View>

        <View className="mb-4">
          <Text className="text-sm text-[#6B7280] dark:text-[#9090a8] mb-2">Last name</Text>
          <Controller
            control={control}
            name="lastName"
            render={({ field: { onChange, onBlur, value } }) => (
              <TextInput
                className="bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-xl px-4 py-3.5 text-[#0F1018] dark:text-white text-base"
                placeholder="Last name"
                placeholderTextColor="#606078"
                value={value}
                onChangeText={onChange}
                onBlur={onBlur}
              />
            )}
          />
          {errors.lastName ? (
            <Text className="text-[#ff6b6b] text-xs mt-1">{errors.lastName.message}</Text>
          ) : null}
        </View>

        <View className="mb-6">
          <Text className="text-sm text-[#6B7280] dark:text-[#9090a8] mb-2">Email</Text>
          <Controller
            control={control}
            name="email"
            render={({ field: { onChange, onBlur, value } }) => (
              <TextInput
                className="bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-xl px-4 py-3.5 text-[#0F1018] dark:text-white text-base"
                placeholder="Email"
                placeholderTextColor="#606078"
                value={value}
                onChangeText={onChange}
                onBlur={onBlur}
                keyboardType="email-address"
                autoCapitalize="none"
              />
            )}
          />
          {errors.email ? (
            <Text className="text-[#ff6b6b] text-xs mt-1">{errors.email.message}</Text>
          ) : null}
        </View>

        <TouchableOpacity
          className={` + "`" + `rounded-xl py-4 items-center mb-3 ${saving ? "bg-[#6c5ce7]/50" : "bg-[#6c5ce7]"}` + "`" + `}
          onPress={handleSubmit(onSave)}
          disabled={saving}
          activeOpacity={0.8}
        >
          {saving ? (
            <ActivityIndicator color="#fff" />
          ) : (
            <Text className="text-white font-semibold text-base">Save changes</Text>
          )}
        </TouchableOpacity>

        <TouchableOpacity
          className="rounded-xl py-4 items-center border border-[#E5E7EB] dark:border-[#2a2a3a]"
          onPress={onCancel}
          activeOpacity={0.8}
        >
          <Text className="text-[#6B7280] dark:text-[#9090a8] font-semibold text-base">Cancel</Text>
        </TouchableOpacity>
      </ScrollView>
    );
  }

  return (
    <ScrollView className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]" contentContainerClassName="px-6 pt-16 pb-28">
      <View className="items-center mb-8 mt-4">
        <TouchableOpacity onPress={onChangeAvatar} activeOpacity={0.85} className="mb-4">
          <View className="w-24 h-24 rounded-full bg-[#6c5ce7] items-center justify-center overflow-hidden">
            {user?.avatar ? (
              <Image source={{ uri: user.avatar }} style={{ width: "100%", height: "100%" }} contentFit="cover" />
            ) : (
              <Text className="text-3xl font-bold text-white">{initials}</Text>
            )}
          </View>
          <View className="absolute bottom-0 right-0 w-8 h-8 rounded-full bg-[#6c5ce7] border-2 border-[#F4F4F6] dark:border-[#0a0a0f] items-center justify-center">
            {uploadingAvatar ? (
              <ActivityIndicator color="#fff" size="small" />
            ) : (
              <Ionicons name="camera" size={15} color="#fff" />
            )}
          </View>
        </TouchableOpacity>
        <Text className="text-xl font-bold text-[#0F1018] dark:text-white">{displayName.trim()}</Text>
        <Text className="text-sm text-[#6B7280] dark:text-[#9090a8] mt-1">{user?.email}</Text>
        <View className="bg-[#6c5ce7]/20 px-3 py-1 rounded-full mt-2">
          <Text className="text-[#6c5ce7] text-xs font-medium capitalize">
            {user?.role || "user"}
          </Text>
        </View>
      </View>

      <View className="bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-2xl overflow-hidden mb-6">
        <ProfileRow icon="person-outline" label="Full Name" value={displayName.trim()} />
        <View className="h-px bg-[#E5E7EB] dark:bg-[#2a2a3a]" />
        <ProfileRow icon="mail-outline" label="Email" value={user?.email || "—"} />
        <View className="h-px bg-[#E5E7EB] dark:bg-[#2a2a3a]" />
        <ProfileRow icon="shield-outline" label="Role" value={user?.role || "user"} />
      </View>

      <TouchableOpacity
        className="bg-[#6c5ce7] rounded-2xl py-4 items-center mb-3"
        onPress={() => setEditing(true)}
        activeOpacity={0.8}
      >
        <Text className="text-white font-semibold text-base">Edit Profile</Text>
      </TouchableOpacity>

      <TouchableOpacity
        className="bg-white dark:bg-[#111118] border border-[#E5E7EB] dark:border-[#2a2a3a] rounded-2xl py-4 items-center mb-3 flex-row justify-center"
        onPress={() => router.push("/change-password")}
        activeOpacity={0.8}
      >
        <Ionicons name="lock-closed-outline" size={18} color="#6c5ce7" />
        <Text className="text-[#0F1018] dark:text-white font-semibold text-base ml-2">Change Password</Text>
      </TouchableOpacity>

      <TouchableOpacity
        className="bg-[#ff6b6b]/10 border border-[#ff6b6b]/30 rounded-2xl py-4 items-center"
        onPress={logout}
        activeOpacity={0.8}
      >
        <Text className="text-[#ff6b6b] font-semibold text-base">Sign out</Text>
      </TouchableOpacity>
    </ScrollView>
  );
}
`
}

func expoSettingsScreen() string {
	return `import { useState } from "react";
import {
  View,
  Text,
  TouchableOpacity,
  SectionList,
  Switch,
  Alert,
} from "react-native";
import { Ionicons } from "@expo/vector-icons";
import * as Haptics from "expo-haptics";
import { useAuth } from "@/lib/auth";
import { useTheme } from "@/lib/theme";

interface SettingItem {
  id: string;
  title: string;
  icon: string;
  type: "toggle" | "select" | "action" | "info";
  value?: string;
  danger?: boolean;
}

interface SettingSection {
  title: string;
  data: SettingItem[];
}

export default function SettingsScreen() {
  const { logout } = useAuth();
  const { scheme, setMode } = useTheme();
  const [notifications, setNotifications] = useState(true);
  const [language, setLanguage] = useState("English");

  const handleToggle = (id: string, value: boolean) => {
    Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Light);
    // Persisted theme switch — flips the whole app between light and dark.
    if (id === "dark_mode") setMode(value ? "dark" : "light");
    if (id === "notifications") setNotifications(value);
  };

  const handleLanguage = () => {
    Alert.alert("Language", "Select your preferred language", [
      { text: "English", onPress: () => setLanguage("English") },
      { text: "Spanish", onPress: () => setLanguage("Spanish") },
      { text: "French", onPress: () => setLanguage("French") },
      { text: "Cancel", style: "cancel" },
    ]);
  };

  const handleClearCache = () => {
    Haptics.notificationAsync(Haptics.NotificationFeedbackType.Warning);
    Alert.alert(
      "Clear Cache",
      "This will clear all cached data. Continue?",
      [
        { text: "Cancel", style: "cancel" },
        {
          text: "Clear",
          style: "destructive",
          onPress: () => {
            Haptics.notificationAsync(Haptics.NotificationFeedbackType.Success);
          },
        },
      ]
    );
  };

  const handleLogout = () => {
    Haptics.notificationAsync(Haptics.NotificationFeedbackType.Warning);
    Alert.alert("Sign Out", "Are you sure you want to sign out?", [
      { text: "Cancel", style: "cancel" },
      {
        text: "Sign Out",
        style: "destructive",
        onPress: () => logout(),
      },
    ]);
  };

  const sections: SettingSection[] = [
    {
      title: "Preferences",
      data: [
        { id: "dark_mode", title: "Dark Mode", icon: "moon-outline", type: "toggle" },
        { id: "notifications", title: "Notifications", icon: "notifications-outline", type: "toggle" },
        { id: "language", title: "Language", icon: "language-outline", type: "select", value: language },
      ],
    },
    {
      title: "Data",
      data: [
        { id: "clear_cache", title: "Clear Cache", icon: "trash-outline", type: "action" },
      ],
    },
    {
      title: "About",
      data: [
        { id: "version", title: "App Version", icon: "information-circle-outline", type: "info", value: "1.0.0" },
        { id: "build", title: "Build", icon: "code-outline", type: "info", value: "1" },
      ],
    },
    {
      title: "",
      data: [
        { id: "logout", title: "Sign Out", icon: "log-out-outline", type: "action", danger: true },
      ],
    },
  ];

  const renderItem = ({ item }: { item: SettingItem }) => {
    const onPress = () => {
      if (item.id === "language") handleLanguage();
      if (item.id === "clear_cache") handleClearCache();
      if (item.id === "logout") handleLogout();
    };

    return (
      <TouchableOpacity
        className="flex-row items-center bg-white dark:bg-[#111118] px-5 py-4"
        activeOpacity={item.type === "info" ? 1 : 0.7}
        onPress={item.type !== "info" && item.type !== "toggle" ? onPress : undefined}
        disabled={item.type === "info" || item.type === "toggle"}
      >
        <Ionicons
          name={item.icon as any}
          size={20}
          color={item.danger ? "#ff6b6b" : "#9090a8"}
        />
        <Text
          className={` + "`" + `text-sm ml-3 flex-1 ${item.danger ? "text-[#ff6b6b] font-semibold" : "text-[#0F1018] dark:text-white"}` + "`" + `}
        >
          {item.title}
        </Text>

        {item.type === "toggle" ? (
          <Switch
            value={item.id === "dark_mode" ? scheme === "dark" : notifications}
            onValueChange={(val) => handleToggle(item.id, val)}
            trackColor={{ false: scheme === "dark" ? "#2a2a3a" : "#D1D5DB", true: "#6c5ce7" }}
            thumbColor="#ffffff"
          />
        ) : null}

        {item.type === "select" ? (
          <View className="flex-row items-center">
            <Text className="text-sm text-[#6B7280] dark:text-[#9090a8] mr-2">{item.value}</Text>
            <Ionicons name="chevron-forward" size={16} color="#606078" />
          </View>
        ) : null}

        {item.type === "info" ? (
          <Text className="text-sm text-[#9CA3AF] dark:text-[#606078]">{item.value}</Text>
        ) : null}

        {item.type === "action" && !item.danger ? (
          <Ionicons name="chevron-forward" size={16} color="#606078" />
        ) : null}
      </TouchableOpacity>
    );
  };

  const renderSectionHeader = ({ section }: { section: SettingSection }) => {
    if (!section.title) return <View className="h-6" />;
    return (
      <View className="px-5 pt-6 pb-2 bg-[#F4F4F6] dark:bg-[#0a0a0f]">
        <Text className="text-xs font-semibold text-[#9CA3AF] dark:text-[#606078] uppercase tracking-wider">
          {section.title}
        </Text>
      </View>
    );
  };

  const renderSeparator = () => <View className="h-px bg-[#E5E7EB] dark:bg-[#2a2a3a] ml-14" />;

  return (
    <SectionList
      className="flex-1 bg-[#F4F4F6] dark:bg-[#0a0a0f]"
      sections={sections}
      keyExtractor={(item) => item.id}
      renderItem={renderItem}
      renderSectionHeader={renderSectionHeader}
      ItemSeparatorComponent={renderSeparator}
      stickySectionHeadersEnabled={false}
      contentContainerClassName="pt-14 pb-28"
    />
  );
}
`
}

func expoAPIClient() string {
	return `import * as SecureStore from "expo-secure-store";
import Constants from "expo-constants";
import { Platform } from "react-native";

// Resolve the API base URL so it works on a simulator AND a physical
// device without hand-editing IPs:
//   1. EXPO_PUBLIC_API_URL wins if set (e.g. a deployed backend).
//   2. Otherwise derive the dev machine's LAN IP from the host the phone
//      already used to reach Metro (Constants ... hostUri). A real device
//      can't reach "localhost"/"10.0.2.2" — those point at the device
//      itself — but it CAN reach whatever IP Metro is served on.
//   3. Fall back to the platform loopback for web / edge cases.
const API_PORT = 8080;

function resolveApiUrl(): string {
  const explicit = process.env.EXPO_PUBLIC_API_URL;
  if (explicit) return explicit.replace(/\/$/, "") + "/api";

  const hostUri =
    Constants.expoConfig?.hostUri ||
    (Constants.expoGoConfig as any)?.debuggerHost ||
    (Constants.manifest2 as any)?.extra?.expoGo?.debuggerHost;
  const host = hostUri ? String(hostUri).split(":")[0] : undefined;
  if (host && host !== "localhost" && host !== "127.0.0.1") {
    return "http://" + host + ":" + API_PORT + "/api";
  }

  return Platform.select({
    android: "http://10.0.2.2:" + API_PORT + "/api",
    default: "http://localhost:" + API_PORT + "/api",
  }) as string;
}

const API_URL = resolveApiUrl();

// Fail fast instead of letting fetch hang for minutes when the API is
// unreachable — that hang is what leaves the splash screen stuck.
const REQUEST_TIMEOUT_MS = 15000;

async function fetchWithTimeout(url: string, init: RequestInit): Promise<Response> {
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS);
  try {
    return await fetch(url, { ...init, signal: controller.signal });
  } finally {
    clearTimeout(timer);
  }
}

export { API_URL };

interface RequestOptions {
  method?: string;
  body?: any;
  headers?: Record<string, string>;
}

// UUIDv4-shaped string for the Idempotency-Key header. Math.random is fine
// here — collision risk for per-mutation keys with 122 bits of entropy is
// effectively zero, and we don't need cryptographic strength for dedupe.
function randomKey(): string {
  return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    const v = c === "x" ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

class ApiClient {
  private async getToken(): Promise<string | null> {
    return SecureStore.getItemAsync("access_token");
  }

  // Public: lets the auth provider skip the /auth/me boot request entirely
  // when there's no session — no token means no network call, so a fresh
  // install dismisses the splash instantly instead of waiting on a fetch.
  async hasToken(): Promise<boolean> {
    return !!(await SecureStore.getItemAsync("access_token"));
  }

  private async getRefreshToken(): Promise<string | null> {
    return SecureStore.getItemAsync("refresh_token");
  }

  async setTokens(accessToken?: string, refreshToken?: string) {
    // Guard against a shape mismatch: SecureStore throws an opaque
    // "Values must be strings" error on undefined, so surface a clear one.
    if (!accessToken || !refreshToken) {
      throw new Error("Auth response did not include tokens");
    }
    await SecureStore.setItemAsync("access_token", accessToken);
    await SecureStore.setItemAsync("refresh_token", refreshToken);
  }

  async clearTokens() {
    await SecureStore.deleteItemAsync("access_token");
    await SecureStore.deleteItemAsync("refresh_token");
  }

  private async request(endpoint: string, options: RequestOptions = {}) {
    const token = await this.getToken();
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      ...options.headers,
    };
    if (token) {
      headers["Authorization"] = ` + "`" + `Bearer ${token}` + "`" + `;
    }

    // Stable idempotency key for unsafe methods so the 401-refresh retry
    // below replays the exact same request and the server can dedupe.
    const method = options.method || "GET";
    const unsafe = method === "POST" || method === "PUT" || method === "PATCH" || method === "DELETE";
    if (unsafe && !headers["Idempotency-Key"]) {
      headers["Idempotency-Key"] = randomKey();
    }

    let res = await fetchWithTimeout(` + "`" + `${API_URL}${endpoint}` + "`" + `, {
      method,
      headers,
      body: options.body ? JSON.stringify(options.body) : undefined,
    });

    // Skip refresh on the auth endpoints themselves — a wrong password
    // 401-ing into a refresh attempt would loop and wipe the session.
    const isAuthEndpoint =
      endpoint.includes("/auth/login") ||
      endpoint.includes("/auth/register") ||
      endpoint.includes("/auth/refresh");

    // Try refresh if unauthorized
    if (res.status === 401 && !isAuthEndpoint) {
      const refreshToken = await this.getRefreshToken();
      if (refreshToken) {
        const refreshRes = await fetchWithTimeout(` + "`" + `${API_URL}/auth/refresh` + "`" + `, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ refresh_token: refreshToken }),
        });

        if (refreshRes.ok) {
          const data = await refreshRes.json();
          await this.setTokens(data.data.tokens.access_token, data.data.tokens.refresh_token);

          headers["Authorization"] = ` + "`" + `Bearer ${data.data.access_token}` + "`" + `;
          res = await fetchWithTimeout(` + "`" + `${API_URL}${endpoint}` + "`" + `, {
            method: options.method || "GET",
            headers,
            body: options.body ? JSON.stringify(options.body) : undefined,
          });
        } else {
          await this.clearTokens();
          throw new Error("Session expired");
        }
      }
    }

    const json = await res.json();
    if (!res.ok) {
      throw new Error(json.error?.message || "Request failed");
    }
    return json;
  }

  get(endpoint: string) {
    return this.request(endpoint, { method: "GET" });
  }

  post(endpoint: string, body: any) {
    return this.request(endpoint, { method: "POST", body });
  }

  put(endpoint: string, body: any) {
    return this.request(endpoint, { method: "PUT", body });
  }

  delete(endpoint: string) {
    return this.request(endpoint, { method: "DELETE" });
  }
}

export const api = new ApiClient();
`
}

func expoAuthProvider() string {
	return `import { createContext, useContext, useState, useEffect, useCallback } from "react";
import * as WebBrowser from "expo-web-browser";
import { api, API_URL } from "./api";

interface User {
  id: string;
  first_name: string;
  last_name: string;
  name?: string;
  email: string;
  role: string;
}

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  loginWithGoogle: () => Promise<void>;
  register: (firstName: string, lastName: string, email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType>({
  user: null,
  isAuthenticated: false,
  isLoading: true,
  login: async () => {},
  loginWithGoogle: async () => {},
  register: async () => {},
  logout: async () => {},
  refreshUser: async () => {},
});

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadUser();
  }, []);

  const loadUser = async () => {
    try {
      // No stored session → don't touch the network. This is what keeps a
      // fresh install from sitting on the splash screen while a doomed
      // /auth/me request waits out its timeout.
      if (!(await api.hasToken())) {
        setUser(null);
        return;
      }
      const res = await api.get("/auth/me");
      setUser(res.data);
    } catch {
      setUser(null);
    } finally {
      setIsLoading(false);
    }
  };

  const login = useCallback(async (email: string, password: string) => {
    const res = await api.post("/auth/login", { email, password });
    await api.setTokens(res.data.tokens.access_token, res.data.tokens.refresh_token);
    setUser(res.data.user);
  }, []);

  const loginWithGoogle = useCallback(async () => {
    const callbackUrl = "myapp://callback";
    const result = await WebBrowser.openAuthSessionAsync(
      ` + "`" + `${API_URL}/auth/oauth/google?redirect_uri=${encodeURIComponent(callbackUrl)}` + "`" + `,
      callbackUrl
    );

    if (result.type === "success" && result.url) {
      const url = new URL(result.url);
      const accessToken = url.searchParams.get("access_token");
      const refreshToken = url.searchParams.get("refresh_token");

      if (accessToken && refreshToken) {
        await api.setTokens(accessToken, refreshToken);
        await loadUser();
      } else {
        throw new Error("OAuth callback missing tokens");
      }
    } else if (result.type === "cancel") {
      throw new Error("Login cancelled");
    }
  }, []);

  const register = useCallback(async (firstName: string, lastName: string, email: string, password: string) => {
    const res = await api.post("/auth/register", { first_name: firstName, last_name: lastName, email, password });
    await api.setTokens(res.data.tokens.access_token, res.data.tokens.refresh_token);
    setUser(res.data.user);
  }, []);

  const logout = useCallback(async () => {
    try {
      await api.post("/auth/logout", {});
    } catch {
      // Ignore errors on logout
    }
    await api.clearTokens();
    setUser(null);
  }, []);

  // Re-fetch the current user (e.g. after a profile/avatar update).
  const refreshUser = useCallback(async () => {
    try {
      const res = await api.get("/auth/me");
      setUser(res.data);
    } catch {
      // Keep the current user on a transient failure.
    }
  }, []);

  return (
    <AuthContext value={{ user, isAuthenticated: !!user, isLoading, login, loginWithGoogle, register, logout, refreshUser }}>
      {children}
    </AuthContext>
  );
}

export const useAuth = () => useContext(AuthContext);
`
}

func expoThemeProvider() string {
	return `import { createContext, useContext, useEffect, useState } from "react";
import { useColorScheme, colorScheme as nwColorScheme } from "nativewind";
import * as SecureStore from "expo-secure-store";

export type ThemeMode = "light" | "dark" | "system";

const STORAGE_KEY = "theme_mode";

// Default to light before the first paint — even on a device whose OS is
// in dark mode — so the app opens light unless the user saved otherwise.
// The provider then reads the persisted choice and overrides if needed.
nwColorScheme.set("light");

// JS-side colours for the things className/dark: variants can't reach:
// gradients, the faint auth grid, blur tint, status bar, tab bar chrome.
export interface Palette {
  scheme: "light" | "dark";
  statusBar: "light" | "dark";
  headerGradient: [string, string, string];
  logoShadow: string;
  gridLine: string;
  gridOpacity: number;
  inputIcon: string;
  placeholder: string;
  tabBar: string;
  tabBarBorder: string;
  tabInactive: string;
  blurTint: "light" | "dark";
  refresh: string;
}

const LIGHT: Palette = {
  scheme: "light",
  statusBar: "dark",
  headerGradient: ["#EFEBFF", "#F6F4FF", "#FFFFFF"],
  logoShadow: "#6c5ce7",
  gridLine: "#0f1018",
  gridOpacity: 0.04,
  inputIcon: "#9CA3AF",
  placeholder: "#9CA3AF",
  tabBar: "#FFFFFF",
  tabBarBorder: "#E5E7EB",
  tabInactive: "#9CA3AF",
  blurTint: "light",
  refresh: "#6c5ce7",
};

const DARK: Palette = {
  scheme: "dark",
  statusBar: "light",
  headerGradient: ["#1c1830", "#15121f", "#111118"],
  logoShadow: "#6c5ce7",
  gridLine: "#e8e8f0",
  gridOpacity: 0.06,
  inputIcon: "#606078",
  placeholder: "#606078",
  tabBar: "#15151d",
  tabBarBorder: "#22222e",
  tabInactive: "#606078",
  blurTint: "dark",
  refresh: "#6c5ce7",
};

interface ThemeContextType {
  mode: ThemeMode;
  scheme: "light" | "dark";
  palette: Palette;
  setMode: (mode: ThemeMode) => void;
}

const ThemeContext = createContext<ThemeContextType>({
  mode: "light",
  scheme: "light",
  palette: LIGHT,
  setMode: () => {},
});

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const { colorScheme, setColorScheme } = useColorScheme();
  const [mode, setModeState] = useState<ThemeMode>("light");

  useEffect(() => {
    // Default is light; only override once we've read the saved choice.
    SecureStore.getItemAsync(STORAGE_KEY).then((saved) => {
      const next = (saved as ThemeMode) || "light";
      setModeState(next);
      setColorScheme(next);
    });
  }, []);

  const setMode = (next: ThemeMode) => {
    setModeState(next);
    setColorScheme(next);
    SecureStore.setItemAsync(STORAGE_KEY, next).catch(() => {});
  };

  const scheme: "light" | "dark" = colorScheme === "dark" ? "dark" : "light";
  const palette = scheme === "dark" ? DARK : LIGHT;

  return (
    <ThemeContext value={{ mode, scheme, palette, setMode }}>
      {children}
    </ThemeContext>
  );
}

export const useTheme = () => useContext(ThemeContext);
`
}

func expoQueryClient() string {
	return `import { QueryClient } from "@tanstack/react-query";

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5,
      retry: 1,
    },
  },
});
`
}
