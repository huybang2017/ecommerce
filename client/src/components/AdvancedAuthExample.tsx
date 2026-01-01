// Example: Using React Query hooks directly in components
// This shows advanced usage patterns

"use client";

import { useLogin, useLogout, useProfile } from "@/hooks/useAuth";
import { useState } from "react";

export function AdvancedAuthExample() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  // Queries
  const {
    data: user,
    isLoading: profileLoading,
    isError: profileError,
    refetch: refetchProfile,
  } = useProfile();

  // Mutations
  const loginMutation = useLogin();
  const logoutMutation = useLogout();

  const handleLogin = () => {
    loginMutation.mutate(
      { email, password },
      {
        onSuccess: (data) => {
          console.log("Login successful!", data.user);
          // Navigate or show success
        },
        onError: (error) => {
          console.error("Login failed:", error.response?.data?.error);
          // Show error toast
        },
      }
    );
  };

  const handleLogout = () => {
    logoutMutation.mutate(undefined, {
      onSuccess: () => {
        console.log("Logout successful!");
        window.location.href = "/login";
      },
    });
  };

  return (
    <div className="p-6">
      <h2 className="text-2xl font-bold mb-4">React Query Auth Example</h2>

      {/* Profile Section */}
      <div className="mb-6 p-4 border rounded">
        <h3 className="font-semibold mb-2">Profile Query</h3>
        {profileLoading && <p>Loading profile...</p>}
        {profileError && <p className="text-red-600">Failed to load profile</p>}
        {user && (
          <div>
            <p>Email: {user.email}</p>
            <p>Name: {user.full_name}</p>
            <p>Role: {user.role}</p>
            <button
              onClick={() => refetchProfile()}
              className="mt-2 px-4 py-2 bg-blue-500 text-white rounded"
            >
              Refresh Profile
            </button>
          </div>
        )}
      </div>

      {/* Login Section */}
      <div className="mb-6 p-4 border rounded">
        <h3 className="font-semibold mb-2">Login Mutation</h3>
        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="block w-full mb-2 px-3 py-2 border rounded"
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="block w-full mb-2 px-3 py-2 border rounded"
        />
        <button
          onClick={handleLogin}
          disabled={loginMutation.isPending}
          className="px-4 py-2 bg-green-500 text-white rounded disabled:bg-gray-400"
        >
          {loginMutation.isPending ? "Logging in..." : "Login"}
        </button>

        {loginMutation.isError && (
          <p className="text-red-600 mt-2">
            Error:{" "}
            {loginMutation.error.response?.data?.error ||
              loginMutation.error.message}
          </p>
        )}

        {loginMutation.isSuccess && (
          <p className="text-green-600 mt-2">Login successful!</p>
        )}
      </div>

      {/* Logout Section */}
      <div className="p-4 border rounded">
        <h3 className="font-semibold mb-2">Logout Mutation</h3>
        <button
          onClick={handleLogout}
          disabled={logoutMutation.isPending}
          className="px-4 py-2 bg-red-500 text-white rounded disabled:bg-gray-400"
        >
          {logoutMutation.isPending ? "Logging out..." : "Logout"}
        </button>
      </div>

      {/* Mutation States */}
      <div className="mt-6 p-4 border rounded bg-gray-50">
        <h3 className="font-semibold mb-2">Mutation States (for debugging)</h3>
        <pre className="text-sm">
          {JSON.stringify(
            {
              login: {
                isPending: loginMutation.isPending,
                isError: loginMutation.isError,
                isSuccess: loginMutation.isSuccess,
              },
              logout: {
                isPending: logoutMutation.isPending,
                isError: logoutMutation.isError,
                isSuccess: logoutMutation.isSuccess,
              },
            },
            null,
            2
          )}
        </pre>
      </div>
    </div>
  );
}
