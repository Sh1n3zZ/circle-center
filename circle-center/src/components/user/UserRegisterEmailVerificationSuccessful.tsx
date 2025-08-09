import { CheckCircle, LogIn, ArrowRight } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

interface UserRegisterEmailVerificationSuccessfulProps {
  email: string
  onProceedToLogin: () => void
  className?: string
}

export function UserRegisterEmailVerificationSuccessful({
  email,
  onProceedToLogin,
  className,
}: UserRegisterEmailVerificationSuccessfulProps) {
  const handleLoginRedirect = () => {
    onProceedToLogin()
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
          <CardTitle className="text-xl font-bold text-green-900">
            Verification Successful!
          </CardTitle>
          <CardDescription>
            Your email has been verified and your account is now active
          </CardDescription>
        </CardHeader>
        
        <CardContent className="space-y-6">
          <div className="text-center space-y-2">
            <p className="text-sm text-muted-foreground">
              Verified email address:
            </p>
            <p className="font-medium text-foreground">{email}</p>
          </div>

          <div className="space-y-4">
            <div className="flex items-start gap-3 p-4 bg-green-50 rounded-lg">
              <CheckCircle className="w-5 h-5 text-green-600 mt-0.5 flex-shrink-0" />
              <div className="space-y-1">
                <h4 className="font-medium text-green-900">Account Activated</h4>
                <p className="text-sm text-green-700">
                  Your Circle Center account is now ready to use. You can sign in with your email and password.
                </p>
              </div>
            </div>

            <div className="bg-blue-50 p-4 rounded-lg">
              <h4 className="font-medium text-blue-900 mb-2">What's next?</h4>
              <ul className="text-sm text-blue-700 space-y-1">
                <li className="flex items-center gap-2">
                  <ArrowRight className="w-3 h-3" />
                  Sign in to your account
                </li>
                <li className="flex items-center gap-2">
                  <ArrowRight className="w-3 h-3" />
                  Complete your profile setup
                </li>
                <li className="flex items-center gap-2">
                  <ArrowRight className="w-3 h-3" />
                  Start using Circle Center features
                </li>
              </ul>
            </div>
          </div>

          <div className="space-y-3">
            <Button onClick={handleLoginRedirect} className="w-full">
              <LogIn className="w-4 h-4 mr-2" />
              Sign In Now
            </Button>
            
            <div className="text-center text-xs text-muted-foreground">
              <p>
                Welcome to Circle Center! If you have any questions, feel free to contact our support team.
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
