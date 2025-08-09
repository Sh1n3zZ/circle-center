import { Mail, AlertTriangle, RefreshCw } from "lucide-react"
import { useState } from "react"

import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { userApi } from "@/api/user/user"
import { toast } from "sonner"

interface UserAccountInactiveDialogProps {
  email: string
  onClose: () => void
  onVerificationSent: (email: string) => void
  className?: string
}

export function UserAccountInactiveDialog({
  email,
  onClose,
  onVerificationSent,
  className,
}: UserAccountInactiveDialogProps) {
  const [isResending, setIsResending] = useState(false)

  const handleResendEmail = async () => {
    setIsResending(true)
    try {
      const response = await userApi.resendVerificationEmail({ email })
      
      if (response.data.email_sent) {
        toast.success("Verification email sent successfully!")
        onVerificationSent(email)
      } else {
        const errorMessage = response.data.email_error || "Failed to send verification email"
        toast.error(errorMessage)
      }
    } catch (error: any) {
      const errorMessage = error.response?.data?.message || "Failed to resend verification email"
      toast.error(errorMessage)
    } finally {
      setIsResending(false)
    }
  }

  return (
    <div className={`fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50 ${className || ""}`}>
      <Card className="w-full max-w-md mx-auto">
        <CardHeader className="space-y-1 text-center">
          <div className="flex justify-center mb-4">
            <div className="flex items-center justify-center w-16 h-16 bg-amber-100 rounded-full">
              <AlertTriangle className="w-8 h-8 text-amber-600" />
            </div>
          </div>
          <CardTitle className="text-xl font-bold">Account Not Verified</CardTitle>
          <CardDescription>
            Your account needs email verification before you can sign in
          </CardDescription>
        </CardHeader>
        
        <CardContent className="space-y-4">
          <div className="text-center space-y-2">
            <p className="text-sm text-muted-foreground">
              We need to verify your email address:
            </p>
            <p className="font-medium text-foreground">{email}</p>
          </div>

          <div className="flex items-start gap-3 p-4 bg-blue-50 rounded-lg">
            <Mail className="w-5 h-5 text-blue-600 mt-0.5 flex-shrink-0" />
            <div className="space-y-1">
              <h4 className="font-medium text-blue-900">Verification Required</h4>
              <p className="text-sm text-blue-700">
                We'll send a verification link to your email. Click the link to activate your account.
              </p>
            </div>
          </div>
        </CardContent>

        <CardFooter className="flex flex-col space-y-3">
          <Button 
            onClick={handleResendEmail}
            disabled={isResending}
            className="w-full"
          >
            {isResending ? (
              <>
                <RefreshCw className="w-4 h-4 mr-2 animate-spin" />
                Sending Verification Email...
              </>
            ) : (
              <>
                <Mail className="w-4 h-4 mr-2" />
                Send Verification Email
              </>
            )}
          </Button>
          
          <Button 
            variant="outline" 
            onClick={onClose}
            disabled={isResending}
            className="w-full"
          >
            Cancel
          </Button>
        </CardFooter>
      </Card>
    </div>
  )
}
