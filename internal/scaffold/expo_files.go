package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeExpoFiles(root string, opts Options) error {
	expoRoot := filepath.Join(root, "apps", "expo")

	files := map[string]string{
		filepath.Join(expoRoot, "package.json"):                  expoPackageJSON(opts),
		filepath.Join(expoRoot, "app.json"):                      expoAppJSON(opts),
		filepath.Join(expoRoot, "tsconfig.json"):                 expoTSConfig(),
		filepath.Join(expoRoot, "tailwind.config.js"):            expoTailwindConfig(),
		filepath.Join(expoRoot, "metro.config.js"):               expoMetroConfig(),
		filepath.Join(expoRoot, "global.css"):                    expoGlobalCSS(),
		filepath.Join(expoRoot, "nativewind-env.d.ts"):           expoNativewindEnv(),
		filepath.Join(expoRoot, "app", "_layout.tsx"):            expoRootLayout(),
		filepath.Join(expoRoot, "app", "(auth)", "_layout.tsx"):  expoAuthLayout(),
		filepath.Join(expoRoot, "app", "(auth)", "login.tsx"):    expoLoginScreen(),
		filepath.Join(expoRoot, "app", "(auth)", "register.tsx"): expoRegisterScreen(),
		filepath.Join(expoRoot, "app", "(tabs)", "_layout.tsx"):  expoTabsLayout(),
		filepath.Join(expoRoot, "app", "(tabs)", "index.tsx"):    expoDashboardScreen(),
		filepath.Join(expoRoot, "app", "(tabs)", "profile.tsx"):  expoProfileScreen(),
		filepath.Join(expoRoot, "lib", "api.ts"):                 expoAPIClient(),
		filepath.Join(expoRoot, "lib", "auth.tsx"):               expoAuthProvider(),
		filepath.Join(expoRoot, "lib", "query-client.ts"):        expoQueryClient(),
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
    "react": "19.1.0",
    "react-native": "0.81.0",
    "react-native-safe-area-context": "~5.6.0",
    "react-native-screens": "~4.16.0",
    "@react-navigation/bottom-tabs": "^7.0.0",
    "@expo/vector-icons": "^15.0.0",
    "nativewind": "^4.2.0",
    "@tanstack/react-query": "^5.0.0",
    "react-native-reanimated": "~4.1.0",
    "react-native-gesture-handler": "~2.28.0"
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
import { useAuth } from "@/lib/auth";

export default function LoginScreen() {
  const { login } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleLogin = async () => {
    if (!email || !password) {
      setError("Please fill in all fields");
      return;
    }
    setError("");
    setLoading(true);
    try {
      await login(email, password);
    } catch (err: any) {
      setError(err.message || "Login failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === "ios" ? "padding" : "height"}
      className="flex-1 bg-bg-primary"
    >
      <View className="flex-1 justify-center px-8">
        <Text className="text-4xl font-bold text-white mb-2">Welcome back</Text>
        <Text className="text-base text-[#9090a8] mb-10">
          Sign in to your account
        </Text>

        {error ? (
          <View className="bg-danger/10 border border-danger/30 rounded-xl p-4 mb-6">
            <Text className="text-danger text-sm">{error}</Text>
          </View>
        ) : null}

        <View className="mb-4">
          <Text className="text-sm text-[#9090a8] mb-2">Email</Text>
          <TextInput
            className="bg-bg-secondary border border-border rounded-xl px-4 py-3.5 text-white text-base"
            placeholder="you@example.com"
            placeholderTextColor="#606078"
            value={email}
            onChangeText={setEmail}
            keyboardType="email-address"
            autoCapitalize="none"
            autoComplete="email"
          />
        </View>

        <View className="mb-6">
          <Text className="text-sm text-[#9090a8] mb-2">Password</Text>
          <TextInput
            className="bg-bg-secondary border border-border rounded-xl px-4 py-3.5 text-white text-base"
            placeholder="••••••••"
            placeholderTextColor="#606078"
            value={password}
            onChangeText={setPassword}
            secureTextEntry
          />
        </View>

        <TouchableOpacity
          className={` + "`" + `rounded-xl py-4 items-center ${loading ? "bg-accent/50" : "bg-accent"}` + "`" + `}
          onPress={handleLogin}
          disabled={loading}
          activeOpacity={0.8}
        >
          {loading ? (
            <ActivityIndicator color="#fff" />
          ) : (
            <Text className="text-white font-semibold text-base">Sign in</Text>
          )}
        </TouchableOpacity>

        <View className="flex-row justify-center mt-6">
          <Text className="text-[#9090a8]">Don't have an account? </Text>
          <Link href="/(auth)/register">
            <Text className="text-accent font-semibold">Sign up</Text>
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
import { useAuth } from "@/lib/auth";

export default function RegisterScreen() {
  const { register } = useAuth();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleRegister = async () => {
    if (!name || !email || !password) {
      setError("Please fill in all fields");
      return;
    }
    if (password.length < 8) {
      setError("Password must be at least 8 characters");
      return;
    }
    setError("");
    setLoading(true);
    try {
      await register(name, email, password);
    } catch (err: any) {
      setError(err.message || "Registration failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <KeyboardAvoidingView
      behavior={Platform.OS === "ios" ? "padding" : "height"}
      className="flex-1 bg-bg-primary"
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
          <View className="bg-danger/10 border border-danger/30 rounded-xl p-4 mb-6">
            <Text className="text-danger text-sm">{error}</Text>
          </View>
        ) : null}

        <View className="mb-4">
          <Text className="text-sm text-[#9090a8] mb-2">Full name</Text>
          <TextInput
            className="bg-bg-secondary border border-border rounded-xl px-4 py-3.5 text-white text-base"
            placeholder="John Doe"
            placeholderTextColor="#606078"
            value={name}
            onChangeText={setName}
            autoComplete="name"
          />
        </View>

        <View className="mb-4">
          <Text className="text-sm text-[#9090a8] mb-2">Email</Text>
          <TextInput
            className="bg-bg-secondary border border-border rounded-xl px-4 py-3.5 text-white text-base"
            placeholder="you@example.com"
            placeholderTextColor="#606078"
            value={email}
            onChangeText={setEmail}
            keyboardType="email-address"
            autoCapitalize="none"
            autoComplete="email"
          />
        </View>

        <View className="mb-6">
          <Text className="text-sm text-[#9090a8] mb-2">Password</Text>
          <TextInput
            className="bg-bg-secondary border border-border rounded-xl px-4 py-3.5 text-white text-base"
            placeholder="Min. 8 characters"
            placeholderTextColor="#606078"
            value={password}
            onChangeText={setPassword}
            secureTextEntry
          />
        </View>

        <TouchableOpacity
          className={` + "`" + `rounded-xl py-4 items-center ${loading ? "bg-accent/50" : "bg-accent"}` + "`" + `}
          onPress={handleRegister}
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
            <Text className="text-accent font-semibold">Sign in</Text>
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
          backgroundColor: "#111118",
          borderTopColor: "#2a2a3a",
        },
        tabBarActiveTintColor: "#6c5ce7",
        tabBarInactiveTintColor: "#606078",
      }}
    >
      <Tabs.Screen
        name="index"
        options={{
          title: "Dashboard",
          tabBarIcon: ({ color, size }) => (
            <Ionicons name="grid-outline" size={size} color={color} />
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
    </Tabs>
  );
}
`
}

func expoDashboardScreen() string {
	return `import { View, Text, ScrollView, RefreshControl } from "react-native";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import { useState, useCallback } from "react";

interface Stats {
  total_users: number;
  active_users: number;
  new_today: number;
}

function StatCard({ title, value, color }: { title: string; value: number; color: string }) {
  return (
    <View className="bg-bg-elevated border border-border rounded-2xl p-5 flex-1 min-w-[140px]">
      <Text className="text-sm text-[#9090a8] mb-1">{title}</Text>
      <Text className={` + "`" + `text-3xl font-bold` + "`" + `} style={{ color }}>
        {value}
      </Text>
    </View>
  );
}

export default function DashboardScreen() {
  const { user } = useAuth();
  const [refreshing, setRefreshing] = useState(false);

  const { data: stats, refetch } = useQuery<Stats>({
    queryKey: ["dashboard-stats"],
    queryFn: async () => {
      const res = await api.get("/users?page=1&page_size=1");
      return {
        total_users: res.data?.meta?.total || 0,
        active_users: res.data?.meta?.total || 0,
        new_today: 0,
      };
    },
  });

  const onRefresh = useCallback(async () => {
    setRefreshing(true);
    await refetch();
    setRefreshing(false);
  }, [refetch]);

  return (
    <ScrollView
      className="flex-1 bg-bg-primary"
      contentContainerClassName="p-6"
      refreshControl={
        <RefreshControl refreshing={refreshing} onRefresh={onRefresh} tintColor="#6c5ce7" />
      }
    >
      <Text className="text-2xl font-bold text-white mb-1">
        Welcome back, {user?.name?.split(" ")[0] || "User"}
      </Text>
      <Text className="text-base text-[#9090a8] mb-8">
        Here's an overview of your application.
      </Text>

      <View className="flex-row gap-4 mb-6 flex-wrap">
        <StatCard title="Total Users" value={stats?.total_users || 0} color="#6c5ce7" />
        <StatCard title="Active" value={stats?.active_users || 0} color="#00b894" />
      </View>

      <View className="bg-bg-secondary border border-border rounded-2xl p-6">
        <Text className="text-lg font-semibold text-white mb-3">Quick Start</Text>
        <Text className="text-sm text-[#9090a8] leading-6">
          Your Grit mobile app is connected to the API. Edit this dashboard in{"\n"}
          apps/expo/app/(tabs)/index.tsx
        </Text>
      </View>
    </ScrollView>
  );
}
`
}

func expoProfileScreen() string {
	return `import { View, Text, TouchableOpacity, ScrollView } from "react-native";
import { useAuth } from "@/lib/auth";
import { Ionicons } from "@expo/vector-icons";

export default function ProfileScreen() {
  const { user, logout } = useAuth();

  return (
    <ScrollView className="flex-1 bg-bg-primary" contentContainerClassName="p-6">
      <View className="items-center mb-8 mt-4">
        <View className="w-20 h-20 rounded-full bg-accent items-center justify-center mb-4">
          <Text className="text-3xl font-bold text-white">
            {user?.name?.charAt(0)?.toUpperCase() || "U"}
          </Text>
        </View>
        <Text className="text-xl font-bold text-white">{user?.name || "User"}</Text>
        <Text className="text-sm text-[#9090a8] mt-1">{user?.email}</Text>
        <View className="bg-accent/20 px-3 py-1 rounded-full mt-2">
          <Text className="text-accent text-xs font-medium capitalize">
            {user?.role || "user"}
          </Text>
        </View>
      </View>

      <View className="bg-bg-secondary border border-border rounded-2xl overflow-hidden mb-6">
        <ProfileRow icon="person-outline" label="Full Name" value={user?.name || "—"} />
        <View className="h-px bg-border" />
        <ProfileRow icon="mail-outline" label="Email" value={user?.email || "—"} />
        <View className="h-px bg-border" />
        <ProfileRow icon="shield-outline" label="Role" value={user?.role || "user"} />
      </View>

      <TouchableOpacity
        className="bg-danger/10 border border-danger/30 rounded-2xl py-4 items-center"
        onPress={logout}
        activeOpacity={0.8}
      >
        <Text className="text-danger font-semibold text-base">Sign out</Text>
      </TouchableOpacity>
    </ScrollView>
  );
}

function ProfileRow({ icon, label, value }: { icon: string; label: string; value: string }) {
  return (
    <View className="flex-row items-center px-5 py-4">
      <Ionicons name={icon as any} size={20} color="#9090a8" />
      <Text className="text-sm text-[#9090a8] ml-3 flex-1">{label}</Text>
      <Text className="text-sm text-white">{value}</Text>
    </View>
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
import { api } from "./api";

interface User {
  id: string;
  name: string;
  email: string;
  role: string;
}

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (name: string, email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType>({
  user: null,
  isAuthenticated: false,
  isLoading: true,
  login: async () => {},
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

  const register = useCallback(async (name: string, email: string, password: string) => {
    const res = await api.post("/auth/register", { name, email, password });
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
    <AuthContext value={{ user, isAuthenticated: !!user, isLoading, login, register, logout }}>
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
