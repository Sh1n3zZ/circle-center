import { useState } from "react"
import { useNavigate, useSearchParams } from "react-router-dom"
import { toast } from "sonner"

import { AuthForm } from "@/components/user/UserAccountForm"
import type { LoginFormData, RegisterFormData } from "@/components/user/UserAccountForm"
import { userApi } from "@/api/user/user"

export default function Login() {
  const [isLoading, setIsLoading] = useState(false)
  const navigate = useNavigate()
  const [searchParams, setSearchParams] = useSearchParams()
  const initialTab = (searchParams.get("tab") as "login" | "register") || "login"

  const handleLogin = async (data: LoginFormData) => {
    setIsLoading(true)
    
    try {
      const response = await userApi.login({
        email: data.email,
        password: data.password,
      })
      
      // Store token if provided
      if (response.data.token) {
        localStorage.setItem("token", response.data.token)
      }
      
      toast.success(response.message || "Login successful!")
      navigate("/")
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || "Login failed. Please try again."
      toast.error(errorMessage)
      console.error("Login error:", error)
    } finally {
      setIsLoading(false)
    }
  }

  const handleRegister = async (data: RegisterFormData) => {
    setIsLoading(true)
    
    try {
      const response = await userApi.register({
        username: data.username,
        email: data.email,
        password: data.password,
        display_name: data.displayName,
        phone: data.phone,
      })
      
      toast.success(response.message || "Registration successful!")
      // Stay on the same page after registration, user can switch to login tab
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || "Registration failed. Please try again."
      toast.error(errorMessage)
      console.error("Registration error:", error)
    } finally {
      setIsLoading(false)
    }
  }

  const handleTabChange = (tab: "login" | "register") => {
    // Update URL query parameter when tab changes
    setSearchParams({ tab })
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <AuthForm
        onLogin={handleLogin}
        onRegister={handleRegister}
        isLoading={isLoading}
        className="w-full max-w-md"
        defaultTab={initialTab}
        onTabChange={handleTabChange}
      />
    </div>
  )
}
