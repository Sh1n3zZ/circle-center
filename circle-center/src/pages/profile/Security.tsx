import React from "react";
import { UserProfileSecurityPanel } from "@/components/user/UserProfileSecurityPanel";

const Security: React.FC = () => {
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-4xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Security Settings</h1>
          <p className="text-gray-600 mt-2">
            Manage your account security and privacy settings.
          </p>
        </div>
        
        <UserProfileSecurityPanel />
      </div>
    </div>
  );
};

export default Security;
