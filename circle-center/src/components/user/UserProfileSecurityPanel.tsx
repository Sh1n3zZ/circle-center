import React from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Shield, Lock, Clock, Eye } from "lucide-react";

interface UserProfileSecurityPanelProps {
  /**
   * Additional CSS classes
   */
  className?: string;
}

export const UserProfileSecurityPanel: React.FC<UserProfileSecurityPanelProps> = ({ className }) => {
  return (
    <Card className={className}>
      <CardHeader>
        <CardTitle>Security Settings</CardTitle>
        <CardDescription>Manage your account security and privacy settings.</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="text-center py-12">
          <div className="mx-auto w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mb-6">
            <Shield className="w-12 h-12 text-gray-400" />
          </div>
          
          <h3 className="text-xl font-semibold text-gray-900 mb-2">
            Under Construction
          </h3>
          
          <p className="text-gray-500 mb-6 max-w-md mx-auto">
            Security settings are currently being developed. This section will include password 
            management, two-factor authentication, and other security features.
          </p>
          
          <div className="space-y-3">
            <div className="flex items-center justify-center space-x-2 text-sm text-gray-400">
              <Lock className="w-4 h-4" />
              <span>Password Management</span>
            </div>
            <div className="flex items-center justify-center space-x-2 text-sm text-gray-400">
              <Shield className="w-4 h-4" />
              <span>Two-Factor Authentication</span>
            </div>
            <div className="flex items-center justify-center space-x-2 text-sm text-gray-400">
              <Clock className="w-4 h-4" />
              <span>Login History</span>
            </div>
            <div className="flex items-center justify-center space-x-2 text-sm text-gray-400">
              <Eye className="w-4 h-4" />
              <span>Privacy Settings</span>
            </div>
          </div>
          
          <div className="mt-8">
            <Button variant="outline" disabled>
              Coming Soon
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};
