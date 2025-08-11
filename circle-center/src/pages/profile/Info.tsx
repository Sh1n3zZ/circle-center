import React from "react";
import { UserProfilePanel } from "@/components/user/UserProfilePanel";

const Info: React.FC = () => {
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-4xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Profile Information</h1>
          <p className="text-gray-600 mt-2">
            Manage your personal information and account settings.
          </p>
        </div>
        
        <UserProfilePanel />
      </div>
    </div>
  );
};

export default Info;
