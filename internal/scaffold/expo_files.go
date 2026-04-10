package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeExpoFiles(root string, opts Options) error {
	expoRoot := filepath.Join(root, "apps", "expo")

	files := map[string]string{
		filepath.Join(expoRoot, "package.json"):                     expoPackageJSON(opts),
		filepath.Join(expoRoot, "app.json"):                         expoAppJSON(opts),
		filepath.Join(expoRoot, "tsconfig.json"):                    expoTSConfig(),
		filepath.Join(expoRoot, "tailwind.config.js"):               expoTailwindConfig(),
		filepath.Join(expoRoot, "metro.config.js"):                  expoMetroConfig(),
		filepath.Join(expoRoot, "global.css"):                       expoGlobalCSS(),
		filepath.Join(expoRoot, "nativewind-env.d.ts"):              expoNativewindEnv(),
		filepath.Join(expoRoot, "app", "_layout.tsx"):               expoRootLayout(),
		filepath.Join(expoRoot, "app", "(auth)", "_layout.tsx"):     expoAuthLayout(),
		filepath.Join(expoRoot, "app", "(auth)", "login.tsx"):       expoLoginScreen(),
		filepath.Join(expoRoot, "app", "(auth)", "register.tsx"):    expoRegisterScreen(),
		filepath.Join(expoRoot, "app", "(tabs)", "_layout.tsx"):     expoTabsLayout(),
		filepath.Join(expoRoot, "app", "(tabs)", "index.tsx"):       expoHomeScreen(),
		filepath.Join(expoRoot, "app", "(tabs)", "explore.tsx"):     expoExploreScreen(),
		filepath.Join(expoRoot, "app", "(tabs)", "profile.tsx"):     expoProfileScreen(),
		filepath.Join(expoRoot, "app", "(tabs)", "settings.tsx"):    expoSettingsScreen(),
		filepath.Join(expoRoot, "lib", "api.ts"):                    expoAPIClient(),
		filepath.Join(expoRoot, "lib", "auth.tsx"):                  expoAuthProvider(),
		filepath.Join(expoRoot, "lib", "query-client.ts"):           expoQueryClient(),
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
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
    "expo-image": "~2.3.0",
    "expo-haptics": "~14.1.0",
    "expo-web-browser": "~14.0.0",
    "react": "19.1.0",
    "react-native": "0.81.0",
    "react-native-safe-area-context": "~5.6.0",
    "react-native-screens": "~4.16.0",
    "@react-navigation/bottom-tabs": "^7.0.0",
    "@expo/vector-icons": "^15.0.0",
    "nativewind": "^4.2.0",
    "@tanstack/react-query": "^5.0.0",
    "react-native-reanimated": "~4.1.0",
    "react-native-gesture-handler": "~2.28.0",
    "react-hook-form": "^7.54.0",
    "@hookform/resolvers": "^3.9.0",
    "zod": "^3.23.0"
  },
  "devDependencies": {
    "@babel/core": "^7.24.0",
    "@types/react": "~19.1.0",
    "tailwindcss": "^3.4.0",
    "typescript": "~5.6.0"
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
        "backgroundColor": "#0a0a0f"
      },
      "package": "com.%s.app"
    },
    "plugins": [
      "expo-router",
      "expo-secure-store"
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
import { Stack } from "expo-router";
import { StatusBar } from "expo-status-bar";
import * as SplashScreen from "expo-splash-screen";
import { QueryClientProvider } from "@tanstack/react-query";
import { AuthProvider, useAuth } from "@/lib/auth";
import { queryClient } from "@/lib/query-client";

SplashScreen.preventAutoHideAsync();

function RootNav() {
  const { isAuthenticated, isLoading } = useAuth();

  useEffect(() => {
    if (!isLoading) {
      SplashScreen.hideAsync();
    }
  }, [isLoading]);

  if (isLoading) return null;

  return (
    <>
      <StatusBar style="light" />
      <Stack screenOptions={{ headerShown: false }}>
        {isAuthenticated ? (
          <Stack.Screen name="(tabs)" />
        ) : (
          <Stack.Screen name="(auth)" />
        )}
      </Stack>
    </>
  );
}

export default function RootLayout() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <RootNav />
      </AuthProvider>
    </QueryClientProvider>
  );
}
`
}

func expoAuthLayout() string {
	return `import { Stack } from "expo-router";

export default function AuthLayout() {
  return (
    <Stack
      screenOptions={{
        headerShown: false,
        contentStyle: { backgroundColor: "#0a0a0f" },
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
  TouchableOpacity,
  ActivityIndicator,
  KeyboardAvoidingView,
  Platform,
} from "react-native";
import { Link } from "expo-router";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import * as WebBrowser from "expo-web-browser";
import { Ionicons } from "@expo/vector-icons";
import { useAuth } from "@/lib/auth";

const loginSchema = z.object({
  email: z.string().email("Invalid email"),
  password: z.string().min(6, "Minimum 6 characters"),
});

type LoginForm = z.infer<typeof loginSchema>;

export default function LoginScreen() {
  const { login, loginWithGoogle } = useAuth();
  const [loading, setLoading] = useState(false);
  const [googleLoading, setGoogleLoading] = useState(false);
  const [error, setError] = useState("");

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
      await login(data.email, data.password);
    } catch (err: any) {
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

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === "ios" ? "padding" : "height"}
      className="flex-1 bg-[#0a0a0f]"
    >
      <View className="flex-1 justify-center px-8">
        <Text className="text-4xl font-bold text-white mb-2">Welcome back</Text>
        <Text className="text-base text-[#9090a8] mb-10">
          Sign in to your account
        </Text>

        {error ? (
          <View className="bg-[#ff6b6b]/10 border border-[#ff6b6b]/30 rounded-xl p-4 mb-6">
            <Text className="text-[#ff6b6b] text-sm">{error}</Text>
          </View>
        ) : null}

        <View className="mb-4">
          <Text className="text-sm text-[#9090a8] mb-2">Email</Text>
          <Controller
            control={control}
            name="email"
            render={({ field: { onChange, onBlur, value } }) => (
              <TextInput
                className="bg-[#111118] border border-[#2a2a3a] rounded-xl px-4 py-3.5 text-white text-base"
                placeholder="you@example.com"
                placeholderTextColor="#606078"
                value={value}
                onChangeText={onChange}
                onBlur={onBlur}
                keyboardType="email-address"
                autoCapitalize="none"
                autoComplete="email"
              />
            )}
          />
          {errors.email ? (
            <Text className="text-[#ff6b6b] text-xs mt-1">{errors.email.message}</Text>
          ) : null}
        </View>

        <View className="mb-6">
          <Text className="text-sm text-[#9090a8] mb-2">Password</Text>
          <Controller
            control={control}
            name="password"
            render={({ field: { onChange, onBlur, value } }) => (
              <TextInput
                className="bg-[#111118] border border-[#2a2a3a] rounded-xl px-4 py-3.5 text-white text-base"
                placeholder="••••••••"
                placeholderTextColor="#606078"
                value={value}
                onChangeText={onChange}
                onBlur={onBlur}
                secureTextEntry
              />
            )}
          />
          {errors.password ? (
            <Text className="text-[#ff6b6b] text-xs mt-1">{errors.password.message}</Text>
          ) : null}
        </View>

        <TouchableOpacity
          className={` + "`" + `rounded-xl py-4 items-center ${loading ? "bg-[#6c5ce7]/50" : "bg-[#6c5ce7]"}` + "`" + `}
          onPress={handleSubmit(onSubmit)}
          disabled={loading || googleLoading}
          activeOpacity={0.8}
        >
          {loading ? (
            <ActivityIndicator color="#fff" />
          ) : (
            <Text className="text-white font-semibold text-base">Sign in</Text>
          )}
        </TouchableOpacity>

        <View className="flex-row items-center my-6">
          <View className="flex-1 h-px bg-[#2a2a3a]" />
          <Text className="text-[#606078] mx-4 text-sm">or</Text>
          <View className="flex-1 h-px bg-[#2a2a3a]" />
        </View>

        <TouchableOpacity
          className={` + "`" + `rounded-xl py-4 items-center flex-row justify-center border border-[#2a2a3a] ${googleLoading ? "bg-[#22222e]/50" : "bg-[#22222e]"}` + "`" + `}
          onPress={handleGoogleLogin}
          disabled={loading || googleLoading}
          activeOpacity={0.8}
        >
          {googleLoading ? (
            <ActivityIndicator color="#fff" />
          ) : (
            <>
              <Ionicons name="logo-google" size={20} color="#e8e8f0" />
              <Text className="text-white font-semibold text-base ml-3">
                Sign in with Google
              </Text>
            </>
          )}
        </TouchableOpacity>

        <View className="flex-row justify-center mt-6">
          <Text className="text-[#9090a8]">Don't have an account? </Text>
          <Link href="/(auth)/register">
            <Text className="text-[#6c5ce7] font-semibold">Sign up</Text>
          </Link>
        </View>
      </View>
    </KeyboardAvoidingView>
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
  TouchableOpacity,
  ActivityIndicator,
  KeyboardAvoidingView,
  Platform,
  ScrollView,
} from "react-native";
import { Link } from "expo-router";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useAuth } from "@/lib/auth";

const registerSchema = z.object({
  firstName: z.string().min(1, "Required"),
  lastName: z.string().min(1, "Required"),
  email: z.string().email("Invalid email"),
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
      await registerUser(data.firstName, data.lastName, data.email, data.password);
    } catch (err: any) {
      setError(err.message || "Registration failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === "ios" ? "padding" : "height"}
      className="flex-1 bg-[#0a0a0f]"
    >
      <ScrollView
        contentContainerStyle={{ flexGrow: 1, justifyContent: "center" }}
        className="px-8"
        keyboardShouldPersistTaps="handled"
      >
        <Text className="text-4xl font-bold text-white mb-2">Create account</Text>
        <Text className="text-base text-[#9090a8] mb-10">
          Get started with Grit
        </Text>

        {error ? (
          <View className="bg-[#ff6b6b]/10 border border-[#ff6b6b]/30 rounded-xl p-4 mb-6">
            <Text className="text-[#ff6b6b] text-sm">{error}</Text>
          </View>
        ) : null}

        <View className="flex-row gap-3 mb-4">
          <View className="flex-1">
            <Text className="text-sm text-[#9090a8] mb-2">First name</Text>
            <Controller
              control={control}
              name="firstName"
              render={({ field: { onChange, onBlur, value } }) => (
                <TextInput
                  className="bg-[#111118] border border-[#2a2a3a] rounded-xl px-4 py-3.5 text-white text-base"
                  placeholder="John"
                  placeholderTextColor="#606078"
                  value={value}
                  onChangeText={onChange}
                  onBlur={onBlur}
                  autoComplete="given-name"
                />
              )}
            />
            {errors.firstName ? (
              <Text className="text-[#ff6b6b] text-xs mt-1">{errors.firstName.message}</Text>
            ) : null}
          </View>
          <View className="flex-1">
            <Text className="text-sm text-[#9090a8] mb-2">Last name</Text>
            <Controller
              control={control}
              name="lastName"
              render={({ field: { onChange, onBlur, value } }) => (
                <TextInput
                  className="bg-[#111118] border border-[#2a2a3a] rounded-xl px-4 py-3.5 text-white text-base"
                  placeholder="Doe"
                  placeholderTextColor="#606078"
                  value={value}
                  onChangeText={onChange}
                  onBlur={onBlur}
                  autoComplete="family-name"
                />
              )}
            />
            {errors.lastName ? (
              <Text className="text-[#ff6b6b] text-xs mt-1">{errors.lastName.message}</Text>
            ) : null}
          </View>
        </View>

        <View className="mb-4">
          <Text className="text-sm text-[#9090a8] mb-2">Email</Text>
          <Controller
            control={control}
            name="email"
            render={({ field: { onChange, onBlur, value } }) => (
              <TextInput
                className="bg-[#111118] border border-[#2a2a3a] rounded-xl px-4 py-3.5 text-white text-base"
                placeholder="you@example.com"
                placeholderTextColor="#606078"
                value={value}
                onChangeText={onChange}
                onBlur={onBlur}
                keyboardType="email-address"
                autoCapitalize="none"
                autoComplete="email"
              />
            )}
          />
          {errors.email ? (
            <Text className="text-[#ff6b6b] text-xs mt-1">{errors.email.message}</Text>
          ) : null}
        </View>

        <View className="mb-4">
          <Text className="text-sm text-[#9090a8] mb-2">Password</Text>
          <Controller
            control={control}
            name="password"
            render={({ field: { onChange, onBlur, value } }) => (
              <TextInput
                className="bg-[#111118] border border-[#2a2a3a] rounded-xl px-4 py-3.5 text-white text-base"
                placeholder="Min. 8 characters"
                placeholderTextColor="#606078"
                value={value}
                onChangeText={onChange}
                onBlur={onBlur}
                secureTextEntry
              />
            )}
          />
          {errors.password ? (
            <Text className="text-[#ff6b6b] text-xs mt-1">{errors.password.message}</Text>
          ) : null}
        </View>

        <View className="mb-6">
          <Text className="text-sm text-[#9090a8] mb-2">Confirm password</Text>
          <Controller
            control={control}
            name="confirmPassword"
            render={({ field: { onChange, onBlur, value } }) => (
              <TextInput
                className="bg-[#111118] border border-[#2a2a3a] rounded-xl px-4 py-3.5 text-white text-base"
                placeholder="Repeat password"
                placeholderTextColor="#606078"
                value={value}
                onChangeText={onChange}
                onBlur={onBlur}
                secureTextEntry
              />
            )}
          />
          {errors.confirmPassword ? (
            <Text className="text-[#ff6b6b] text-xs mt-1">{errors.confirmPassword.message}</Text>
          ) : null}
        </View>

        <TouchableOpacity
          className={` + "`" + `rounded-xl py-4 items-center ${loading ? "bg-[#6c5ce7]/50" : "bg-[#6c5ce7]"}` + "`" + `}
          onPress={handleSubmit(onSubmit)}
          disabled={loading}
          activeOpacity={0.8}
        >
          {loading ? (
            <ActivityIndicator color="#fff" />
          ) : (
            <Text className="text-white font-semibold text-base">Create account</Text>
          )}
        </TouchableOpacity>

        <View className="flex-row justify-center mt-6 mb-8">
          <Text className="text-[#9090a8]">Already have an account? </Text>
          <Link href="/(auth)/login">
            <Text className="text-[#6c5ce7] font-semibold">Sign in</Text>
          </Link>
        </View>
      </ScrollView>
    </KeyboardAvoidingView>
  );
}
`
}

func expoTabsLayout() string {
	return `import { Tabs } from "expo-router";
import { Ionicons } from "@expo/vector-icons";

export default function TabsLayout() {
  return (
    <Tabs
      screenOptions={{
        headerStyle: { backgroundColor: "#111118" },
        headerTintColor: "#e8e8f0",
        tabBarStyle: {
          backgroundColor: "#0a0a0f",
          borderTopColor: "#2a2a3a",
          height: 60,
          paddingBottom: 8,
        },
        tabBarActiveTintColor: "#6c5ce7",
        tabBarInactiveTintColor: "#606078",
        tabBarLabelStyle: {
          fontSize: 11,
          fontWeight: "500",
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
          title: "Explore",
          tabBarIcon: ({ color, size }) => (
            <Ionicons name="search-outline" size={size} color={color} />
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
    <View className="bg-[#22222e] border border-[#2a2a3a] rounded-2xl p-4 flex-1 min-w-[140px]">
      <View className="flex-row items-center justify-between mb-2">
        <Ionicons name={icon as any} size={18} color={color} />
      </View>
      <Text className="text-2xl font-bold text-white">{value}</Text>
      <Text className="text-xs text-[#9090a8] mt-1">{title}</Text>
    </View>
  );
}

function RecentItemRow({ item }: { item: RecentItem }) {
  return (
    <View className="flex-row items-center bg-[#22222e] border border-[#2a2a3a] rounded-xl px-4 py-3 mb-2">
      <View className="w-10 h-10 rounded-full bg-[#6c5ce7]/20 items-center justify-center mr-3">
        <Ionicons name={item.icon as any} size={18} color="#6c5ce7" />
      </View>
      <View className="flex-1">
        <Text className="text-sm font-medium text-white">{item.title}</Text>
        <Text className="text-xs text-[#9090a8] mt-0.5">{item.subtitle}</Text>
      </View>
      <Text className="text-xs text-[#606078]">{item.time}</Text>
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
      className="flex-1 bg-[#0a0a0f]"
      contentContainerClassName="p-6"
      refreshControl={
        <RefreshControl refreshing={refreshing} onRefresh={onRefresh} tintColor="#6c5ce7" />
      }
    >
      <Text className="text-2xl font-bold text-white mb-1">
        Welcome back, {firstName}
      </Text>
      <Text className="text-base text-[#9090a8] mb-6">
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

      <Text className="text-lg font-semibold text-white mb-3">Recent Activity</Text>

      {recentItems?.map((item) => (
        <RecentItemRow key={item.id} item={item} />
      ))}

      <View className="bg-[#111118] border border-[#2a2a3a] rounded-2xl p-6 mt-4">
        <Text className="text-lg font-semibold text-white mb-3">Quick Start</Text>
        <Text className="text-sm text-[#9090a8] leading-6">
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
	return `import { useState, useCallback } from "react";
import {
  View,
  Text,
  TextInput,
  ScrollView,
  RefreshControl,
  TouchableOpacity,
} from "react-native";
import { Ionicons } from "@expo/vector-icons";

interface ExploreCategory {
  id: string;
  title: string;
  description: string;
  icon: string;
  color: string;
  count: number;
}

const categories: ExploreCategory[] = [
  { id: "1", title: "Users", description: "Manage user accounts", icon: "people-outline", color: "#6c5ce7", count: 0 },
  { id: "2", title: "Content", description: "Posts, pages, and media", icon: "document-text-outline", color: "#00b894", count: 0 },
  { id: "3", title: "Analytics", description: "Usage and performance", icon: "bar-chart-outline", color: "#74b9ff", count: 0 },
  { id: "4", title: "Notifications", description: "Alerts and messages", icon: "notifications-outline", color: "#fdcb6e", count: 0 },
  { id: "5", title: "Storage", description: "Files and uploads", icon: "cloud-outline", color: "#ff6b6b", count: 0 },
  { id: "6", title: "Integrations", description: "Connected services", icon: "extension-puzzle-outline", color: "#a29bfe", count: 0 },
];

function CategoryCard({ item }: { item: ExploreCategory }) {
  return (
    <TouchableOpacity
      className="bg-[#22222e] border border-[#2a2a3a] rounded-2xl p-5 mb-3"
      activeOpacity={0.7}
    >
      <View className="flex-row items-center">
        <View
          className="w-11 h-11 rounded-xl items-center justify-center mr-4"
          style={{ backgroundColor: item.color + "20" }}
        >
          <Ionicons name={item.icon as any} size={22} color={item.color} />
        </View>
        <View className="flex-1">
          <Text className="text-base font-semibold text-white">{item.title}</Text>
          <Text className="text-xs text-[#9090a8] mt-0.5">{item.description}</Text>
        </View>
        <Ionicons name="chevron-forward" size={18} color="#606078" />
      </View>
    </TouchableOpacity>
  );
}

export default function ExploreScreen() {
  const [search, setSearch] = useState("");
  const [refreshing, setRefreshing] = useState(false);

  const filtered = categories.filter(
    (c) =>
      c.title.toLowerCase().includes(search.toLowerCase()) ||
      c.description.toLowerCase().includes(search.toLowerCase())
  );

  const onRefresh = useCallback(async () => {
    setRefreshing(true);
    // Reload data when API endpoints are available
    setRefreshing(false);
  }, []);

  return (
    <ScrollView
      className="flex-1 bg-[#0a0a0f]"
      contentContainerClassName="p-6"
      refreshControl={
        <RefreshControl refreshing={refreshing} onRefresh={onRefresh} tintColor="#6c5ce7" />
      }
    >
      <Text className="text-2xl font-bold text-white mb-1">Explore</Text>
      <Text className="text-base text-[#9090a8] mb-6">
        Browse and discover features.
      </Text>

      <View className="flex-row items-center bg-[#111118] border border-[#2a2a3a] rounded-xl px-4 mb-6">
        <Ionicons name="search-outline" size={18} color="#606078" />
        <TextInput
          className="flex-1 text-white text-base py-3 ml-3"
          placeholder="Search features..."
          placeholderTextColor="#606078"
          value={search}
          onChangeText={setSearch}
          autoCapitalize="none"
        />
        {search ? (
          <TouchableOpacity onPress={() => setSearch("")}>
            <Ionicons name="close-circle" size={18} color="#606078" />
          </TouchableOpacity>
        ) : null}
      </View>

      {filtered.map((item) => (
        <CategoryCard key={item.id} item={item} />
      ))}

      {filtered.length === 0 ? (
        <View className="items-center py-12">
          <Ionicons name="search-outline" size={48} color="#606078" />
          <Text className="text-[#606078] text-base mt-4">No results found</Text>
        </View>
      ) : null}
    </ScrollView>
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
} from "react-native";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useAuth } from "@/lib/auth";
import { api } from "@/lib/api";
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
      <Text className="text-sm text-[#9090a8] ml-3 flex-1">{label}</Text>
      <Text className="text-sm text-white">{value}</Text>
    </View>
  );
}

export default function ProfileScreen() {
  const { user, logout } = useAuth();
  const [editing, setEditing] = useState(false);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

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
      await api.put("/auth/me", {
        first_name: data.firstName,
        last_name: data.lastName,
        email: data.email,
      });
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
      <ScrollView className="flex-1 bg-[#0a0a0f]" contentContainerClassName="p-6">
        <Text className="text-2xl font-bold text-white mb-6">Edit Profile</Text>

        {error ? (
          <View className="bg-[#ff6b6b]/10 border border-[#ff6b6b]/30 rounded-xl p-4 mb-6">
            <Text className="text-[#ff6b6b] text-sm">{error}</Text>
          </View>
        ) : null}

        <View className="mb-4">
          <Text className="text-sm text-[#9090a8] mb-2">First name</Text>
          <Controller
            control={control}
            name="firstName"
            render={({ field: { onChange, onBlur, value } }) => (
              <TextInput
                className="bg-[#111118] border border-[#2a2a3a] rounded-xl px-4 py-3.5 text-white text-base"
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
          <Text className="text-sm text-[#9090a8] mb-2">Last name</Text>
          <Controller
            control={control}
            name="lastName"
            render={({ field: { onChange, onBlur, value } }) => (
              <TextInput
                className="bg-[#111118] border border-[#2a2a3a] rounded-xl px-4 py-3.5 text-white text-base"
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
          <Text className="text-sm text-[#9090a8] mb-2">Email</Text>
          <Controller
            control={control}
            name="email"
            render={({ field: { onChange, onBlur, value } }) => (
              <TextInput
                className="bg-[#111118] border border-[#2a2a3a] rounded-xl px-4 py-3.5 text-white text-base"
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
          className="rounded-xl py-4 items-center border border-[#2a2a3a]"
          onPress={onCancel}
          activeOpacity={0.8}
        >
          <Text className="text-[#9090a8] font-semibold text-base">Cancel</Text>
        </TouchableOpacity>
      </ScrollView>
    );
  }

  return (
    <ScrollView className="flex-1 bg-[#0a0a0f]" contentContainerClassName="p-6">
      <View className="items-center mb-8 mt-4">
        <View className="w-20 h-20 rounded-full bg-[#6c5ce7] items-center justify-center mb-4">
          <Text className="text-3xl font-bold text-white">{initials}</Text>
        </View>
        <Text className="text-xl font-bold text-white">{displayName.trim()}</Text>
        <Text className="text-sm text-[#9090a8] mt-1">{user?.email}</Text>
        <View className="bg-[#6c5ce7]/20 px-3 py-1 rounded-full mt-2">
          <Text className="text-[#6c5ce7] text-xs font-medium capitalize">
            {user?.role || "user"}
          </Text>
        </View>
      </View>

      <View className="bg-[#111118] border border-[#2a2a3a] rounded-2xl overflow-hidden mb-6">
        <ProfileRow icon="person-outline" label="Full Name" value={displayName.trim()} />
        <View className="h-px bg-[#2a2a3a]" />
        <ProfileRow icon="mail-outline" label="Email" value={user?.email || "—"} />
        <View className="h-px bg-[#2a2a3a]" />
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
  const [darkMode, setDarkMode] = useState(true);
  const [notifications, setNotifications] = useState(true);
  const [language, setLanguage] = useState("English");

  const handleToggle = (id: string, value: boolean) => {
    Haptics.impactAsync(Haptics.ImpactFeedbackStyle.Light);
    if (id === "dark_mode") setDarkMode(value);
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
        className="flex-row items-center bg-[#111118] px-5 py-4"
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
          className={` + "`" + `text-sm ml-3 flex-1 ${item.danger ? "text-[#ff6b6b] font-semibold" : "text-white"}` + "`" + `}
        >
          {item.title}
        </Text>

        {item.type === "toggle" ? (
          <Switch
            value={item.id === "dark_mode" ? darkMode : notifications}
            onValueChange={(val) => handleToggle(item.id, val)}
            trackColor={{ false: "#2a2a3a", true: "#6c5ce7" }}
            thumbColor="#e8e8f0"
          />
        ) : null}

        {item.type === "select" ? (
          <View className="flex-row items-center">
            <Text className="text-sm text-[#9090a8] mr-2">{item.value}</Text>
            <Ionicons name="chevron-forward" size={16} color="#606078" />
          </View>
        ) : null}

        {item.type === "info" ? (
          <Text className="text-sm text-[#606078]">{item.value}</Text>
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
      <View className="px-5 pt-6 pb-2 bg-[#0a0a0f]">
        <Text className="text-xs font-semibold text-[#606078] uppercase tracking-wider">
          {section.title}
        </Text>
      </View>
    );
  };

  const renderSeparator = () => <View className="h-px bg-[#2a2a3a] ml-14" />;

  return (
    <SectionList
      className="flex-1 bg-[#0a0a0f]"
      sections={sections}
      keyExtractor={(item) => item.id}
      renderItem={renderItem}
      renderSectionHeader={renderSectionHeader}
      ItemSeparatorComponent={renderSeparator}
      stickySectionHeadersEnabled={false}
      contentContainerClassName="pb-12"
    />
  );
}
`
}

func expoAPIClient() string {
	return `import * as SecureStore from "expo-secure-store";
import { Platform } from "react-native";

const API_URL = Platform.select({
  android: "http://10.0.2.2:8080/api",
  ios: "http://localhost:8080/api",
  default: "http://localhost:8080/api",
});

export { API_URL };

interface RequestOptions {
  method?: string;
  body?: any;
  headers?: Record<string, string>;
}

class ApiClient {
  private async getToken(): Promise<string | null> {
    return SecureStore.getItemAsync("access_token");
  }

  private async getRefreshToken(): Promise<string | null> {
    return SecureStore.getItemAsync("refresh_token");
  }

  async setTokens(accessToken: string, refreshToken: string) {
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

    let res = await fetch(` + "`" + `${API_URL}${endpoint}` + "`" + `, {
      method: options.method || "GET",
      headers,
      body: options.body ? JSON.stringify(options.body) : undefined,
    });

    // Try refresh if unauthorized
    if (res.status === 401) {
      const refreshToken = await this.getRefreshToken();
      if (refreshToken) {
        const refreshRes = await fetch(` + "`" + `${API_URL}/auth/refresh` + "`" + `, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ refresh_token: refreshToken }),
        });

        if (refreshRes.ok) {
          const data = await refreshRes.json();
          await this.setTokens(data.data.access_token, data.data.refresh_token);

          headers["Authorization"] = ` + "`" + `Bearer ${data.data.access_token}` + "`" + `;
          res = await fetch(` + "`" + `${API_URL}${endpoint}` + "`" + `, {
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
}

const AuthContext = createContext<AuthContextType>({
  user: null,
  isAuthenticated: false,
  isLoading: true,
  login: async () => {},
  loginWithGoogle: async () => {},
  register: async () => {},
  logout: async () => {},
});

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadUser();
  }, []);

  const loadUser = async () => {
    try {
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
    await api.setTokens(res.data.access_token, res.data.refresh_token);
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
    await api.setTokens(res.data.access_token, res.data.refresh_token);
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

  return (
    <AuthContext value={{ user, isAuthenticated: !!user, isLoading, login, loginWithGoogle, register, logout }}>
      {children}
    </AuthContext>
  );
}

export const useAuth = () => useContext(AuthContext);
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
