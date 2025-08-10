import { useState } from "react"
import { useNavigate, useSearchParams } from "react-router-dom"
import { toast } from "sonner"

import { AuthForm } from "@/components/user/UserAccountForm"
import type { LoginFormData, RegisterFormData } from "@/components/user/UserAccountForm"
import { UserRegisterEmailSentSuccessfully } from "@/components/user/UserRegisterEmailSentSuccessfully"
import { UserAccountInactiveDialog } from "@/components/user/UserAccountInactiveDialog"
import { UserRegisterEmailVerification } from "@/components/user/UserRegisterEmailVerification"
import { UserRegisterEmailVerificationSuccessful } from "@/components/user/UserRegisterEmailVerificationSuccessful"
import { userApi } from "@/api/user/user"
import { authHelpers } from "@/api/client"

export default function Login() {
  const [isLoading, setIsLoading] = useState(false)
  const [emailSent, setEmailSent] = useState(false)
  const [registeredEmail, setRegisteredEmail] = useState("")
  const [showInactiveDialog, setShowInactiveDialog] = useState(false)
  const [inactiveEmail, setInactiveEmail] = useState("")
  const [verificationSuccess, setVerificationSuccess] = useState(false)
  const [verificationEmail, setVerificationEmail] = useState("")
  const navigate = useNavigate()
  const [searchParams, setSearchParams] = useSearchParams()
  const initialTab = (searchParams.get("tab") as "login" | "register") || "login"

  // verification
  const isVerificationFlow = searchParams.get("verification") === "true"
  const verificationToken = searchParams.get("token") || ""
  const verificationEmailParam = searchParams.get("email") || ""

  const handleLogin = async (data: LoginFormData) => {
    setIsLoading(true)
    
    try {
      const response = await userApi.login({
        email: data.email,
        password: data.password,
      })
      
      if (response.data.token && response.data.expires_at) {
        const success = await authHelpers.storeAuthData(
          response.data.token,
          response.data.expires_at,
          {
            id: response.data.id,
            username: response.data.username,
            email: response.data.email,
            display_name: response.data.display_name,
            phone: response.data.phone,
            locale: response.data.locale,
            timezone: response.data.timezone,
            avatar_url: response.data.avatar_url,
          }
        )
        
        if (!success) {
          throw new Error("Failed to store authentication data")
        }
      }
      
      toast.success(response.message || "Login successful!")
      navigate("/")
    } catch (error: any) {
      // Check if the error is due to unverified account
      if (error.response?.status === 403 && error.response?.data?.code === "ACCOUNT_NOT_VERIFIED") {
        setInactiveEmail(error.response.data.email || data.email)
        setShowInactiveDialog(true)
      } else {
        const errorMessage = error.response?.data?.message || "Login failed. Please try again."
        toast.error(errorMessage)
      }
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
      
      if (response.data.email_sent) {
        setRegisteredEmail(data.email)
        setEmailSent(true)
        toast.success("Registration successful! Please check your email for verification.")
      } else {
        toast.success(response.message || "Registration successful!")
        if (response.data.email_error) {
          toast.error(`Email sending failed: ${response.data.email_error}`)
        }
      }
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || "Registration failed. Please try again."
      toast.error(errorMessage)
      console.error("Registration error:", error)
    } finally {
      setIsLoading(false)
    }
  }

  const handleTabChange = (tab: "login" | "register") => {
    setSearchParams({ tab })
  }

  const handleCloseInactiveDialog = () => {
    setShowInactiveDialog(false)
    setInactiveEmail("")
  }

  const handleVerificationSent = (email: string) => {
    setShowInactiveDialog(false)
    setRegisteredEmail(email)
    setEmailSent(true)
  }

  const handleVerificationSuccess = () => {
    setVerificationSuccess(true)
    setVerificationEmail(verificationEmailParam)
  }

  const handleVerificationFailed = () => {
    // verification failed, stay on verification page to show error
  }

  const handleProceedToLogin = () => {
    // clear all query parameters and return to login form
    navigate("/login", { replace: true })
    setVerificationSuccess(false)
    setVerificationEmail("")
  }

  if (verificationSuccess) {
    return (
      <UserRegisterEmailVerificationSuccessful
        email={verificationEmail}
        onProceedToLogin={handleProceedToLogin}
      />
    )
  }

  if (isVerificationFlow && verificationToken && verificationEmailParam) {
    return (
      <UserRegisterEmailVerification
        token={verificationToken}
        email={verificationEmailParam}
        onVerificationSuccess={handleVerificationSuccess}
        onVerificationFailed={handleVerificationFailed}
      />
    )
  }

  if (emailSent) {
    return <UserRegisterEmailSentSuccessfully email={registeredEmail} />
  }

  return (
    <>
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
      
      {showInactiveDialog && (
        <UserAccountInactiveDialog
          email={inactiveEmail}
          onClose={handleCloseInactiveDialog}
          onVerificationSent={handleVerificationSent}
        />
      )}
    </>
  )
}
