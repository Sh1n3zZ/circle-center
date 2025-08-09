import { Mail, CheckCircle, AlertCircle, RefreshCw } from "lucide-react"
import { useState } from "react"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { userApi } from "@/api/user/user"
import { toast } from "sonner"

interface EmailSentSuccessfullyProps {
  email: string
  className?: string
}

export function UserRegisterEmailSentSuccessfully({ email, className }: EmailSentSuccessfullyProps) {
  const [isResending, setIsResending] = useState(false)

  const handleResendEmail = async () => {
    setIsResending(true)
    try {
      await userApi.resendVerificationEmail({ email })
      toast.success("Verification email sent successfully!")
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || "Failed to resend verification email"
      toast.error(errorMessage)
    } finally {
      setIsResending(false)
    }
  }

  const handleBackToLogin = () => {
    window.location.href = "/login"
  }

  return (
    <div className={`min-h-screen flex items-center justify-center p-4 ${className || ""}`}>
      <Card className="w-full max-w-md mx-auto">
        <CardHeader className="space-y-1 text-center">
          <div className="flex justify-center mb-4">
            <div className="flex items-center justify-center w-16 h-16 bg-green-100 rounded-full">
              <CheckCircle className="w-8 h-8 text-green-600" />
            </div>
          </div>
          <CardTitle className="text-xl font-bold">Registration Successful!</CardTitle>
          <CardDescription>
            We've sent a verification email to your inbox
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="text-center space-y-2">
            <p className="text-sm text-muted-foreground">
              We've sent a verification email to:
            </p>
            <p className="font-medium text-foreground">{email}</p>
          </div>

          <div className="space-y-4">
            <div className="flex items-start gap-3 p-4 bg-blue-50 rounded-lg">
              <Mail className="w-5 h-5 text-blue-600 mt-0.5 flex-shrink-0" />
              <div className="space-y-1">
                <h4 className="font-medium text-blue-900">Check Your Email</h4>
                <p className="text-sm text-blue-700">
                  Click the verification link in the email to activate your account.
                </p>
              </div>
            </div>

            <div className="flex items-start gap-3 p-4 bg-amber-50 rounded-lg">
              <AlertCircle className="w-5 h-5 text-amber-600 mt-0.5 flex-shrink-0" />
              <div className="space-y-1">
                <h4 className="font-medium text-amber-900">Can't Find the Email?</h4>
                <p className="text-sm text-amber-700">
                  Check your spam folder or junk mail. The email might have been filtered there.
                </p>
              </div>
            </div>
          </div>

          <div className="space-y-3">
            <Button onClick={handleBackToLogin} className="w-full">
              Back to Login
            </Button>
            <Button 
              variant="outline" 
              onClick={handleResendEmail}
              disabled={isResending}
              className="w-full"
            >
              {isResending ? (
                <>
                  <RefreshCw className="w-4 h-4 mr-2 animate-spin" />
                  Sending...
                </>
              ) : (
                <>
                  <Mail className="w-4 h-4 mr-2" />
                  Resend Email
                </>
              )}
            </Button>
          </div>

          <div className="text-center text-xs text-muted-foreground">
            <p>
              Didn't receive the email?{" "}
              <button 
                onClick={handleResendEmail}
                disabled={isResending}
                className="underline underline-offset-4 hover:text-primary disabled:opacity-50"
              >
                Click here to resend
              </button>
            </p>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
