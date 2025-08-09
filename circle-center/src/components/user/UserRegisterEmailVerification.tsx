import { useState, useEffect } from "react"
import { Mail, CheckCircle, XCircle, Loader2 } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { verificationApi } from "@/api/user/verification"
import { toast } from "sonner"

interface UserRegisterEmailVerificationProps {
  token: string
  email: string
  onVerificationSuccess: () => void
  onVerificationFailed: () => void
  className?: string
}

export function UserRegisterEmailVerification({
  token,
  email,
  onVerificationSuccess,
  onVerificationFailed,
  className,
}: UserRegisterEmailVerificationProps) {
  const [isVerifying, setIsVerifying] = useState(false)
  const [verificationStatus, setVerificationStatus] = useState<"pending" | "success" | "failed">("pending")
  const [errorMessage, setErrorMessage] = useState("")

  useEffect(() => {
    handleVerification()
  }, [token, email])

  const handleVerification = async () => {
    if (!token || !email) {
      setVerificationStatus("failed")
      setErrorMessage("Missing verification token or email address")
      toast.error("Invalid verification link")
      return
    }

    setIsVerifying(true)
    setVerificationStatus("pending")

    try {
      const response = await verificationApi.verifyEmail({
        token,
        email,
      })

      if (response.success && response.data.success) {
        setVerificationStatus("success")
        toast.success(response.data.message || "Email verified successfully!")
        onVerificationSuccess()
      } else {
        setVerificationStatus("failed")
        const message = response.data.message || "Email verification failed"
        setErrorMessage(message)
        toast.error(message)
        onVerificationFailed()
      }
    } catch (error: any) {
      setVerificationStatus("failed")
      const message = error.response?.data?.message || "Failed to verify email"
      setErrorMessage(message)
      toast.error(message)
      onVerificationFailed()
    } finally {
      setIsVerifying(false)
    }
  }

  const handleRetry = () => {
    handleVerification()
  }

  const getStatusIcon = () => {
    switch (verificationStatus) {
      case "pending":
        return <Loader2 className="w-8 h-8 text-blue-600 animate-spin" />
      case "success":
        return <CheckCircle className="w-8 h-8 text-green-600" />
      case "failed":
        return <XCircle className="w-8 h-8 text-red-600" />
    }
  }

  const getStatusColor = () => {
    switch (verificationStatus) {
      case "pending":
        return "bg-blue-100"
      case "success":
        return "bg-green-100"
      case "failed":
        return "bg-red-100"
    }
  }

  const getStatusTitle = () => {
    switch (verificationStatus) {
      case "pending":
        return "Verifying Email..."
      case "success":
        return "Email Verified!"
      case "failed":
        return "Verification Failed"
    }
  }

  const getStatusDescription = () => {
    switch (verificationStatus) {
      case "pending":
        return "Please wait while we verify your email address"
      case "success":
        return "Your email has been successfully verified"
      case "failed":
        return errorMessage || "We couldn't verify your email address"
    }
  }

  return (
    <div className={`min-h-screen flex items-center justify-center p-4 ${className || ""}`}>
      <Card className="w-full max-w-md mx-auto">
        <CardHeader className="space-y-1 text-center">
          <div className="flex justify-center mb-4">
            <div className={`flex items-center justify-center w-16 h-16 rounded-full ${getStatusColor()}`}>
              {getStatusIcon()}
            </div>
          </div>
          <CardTitle className="text-xl font-bold">{getStatusTitle()}</CardTitle>
          <CardDescription>{getStatusDescription()}</CardDescription>
        </CardHeader>

        <CardContent className="space-y-6">
          <div className="text-center space-y-2">
            <p className="text-sm text-muted-foreground">
              Email address:
            </p>
            <p className="font-medium text-foreground">{email}</p>
          </div>

          {verificationStatus === "pending" && (
            <div className="flex items-center gap-3 p-4 bg-blue-50 rounded-lg">
              <Mail className="w-5 h-5 text-blue-600 flex-shrink-0" />
              <div className="space-y-1">
                <h4 className="font-medium text-blue-900">Processing Verification</h4>
                <p className="text-sm text-blue-700">
                  We're checking your verification token and activating your account.
                </p>
              </div>
            </div>
          )}

          {verificationStatus === "failed" && (
            <div className="space-y-4">
              <div className="flex items-start gap-3 p-4 bg-red-50 rounded-lg">
                <XCircle className="w-5 h-5 text-red-600 mt-0.5 flex-shrink-0" />
                <div className="space-y-1">
                  <h4 className="font-medium text-red-900">Verification Failed</h4>
                  <p className="text-sm text-red-700">
                    {errorMessage || "The verification link may have expired or is invalid."}
                  </p>
                </div>
              </div>

              <div className="space-y-3">
                <Button onClick={handleRetry} disabled={isVerifying} className="w-full">
                  {isVerifying ? (
                    <>
                      <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                      Retrying...
                    </>
                  ) : (
                    "Try Again"
                  )}
                </Button>
                
                <Button 
                  variant="outline" 
                  onClick={() => window.location.href = "/login"}
                  className="w-full"
                >
                  Back to Login
                </Button>
              </div>
            </div>
          )}

          {verificationStatus === "success" && (
            <div className="space-y-4">
              <div className="flex items-start gap-3 p-4 bg-green-50 rounded-lg">
                <CheckCircle className="w-5 h-5 text-green-600 mt-0.5 flex-shrink-0" />
                <div className="space-y-1">
                  <h4 className="font-medium text-green-900">Success!</h4>
                  <p className="text-sm text-green-700">
                    Your account is now active and you can sign in.
                  </p>
                </div>
              </div>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
