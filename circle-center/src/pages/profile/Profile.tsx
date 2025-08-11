import React, { useState } from "react";
import { UserProfileSideBar, type ProfileTab } from "@/components/user/UserProfileSideBar";
import Info from "./Info";
import Security from "./Security";

const Profile: React.FC = () => {
  const [activeTab, setActiveTab] = useState<ProfileTab>("profile");

  const handleTabChange = (tab: ProfileTab) => {
    setActiveTab(tab);
  };

  const renderContent = () => {
    switch (activeTab) {
      case "profile":
        return <Info />;
      case "security":
        return <Security />;
      default:
        return <Info />;
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="flex">
        {/* Sidebar */}
        <UserProfileSideBar
          activeTab={activeTab}
          onTabChange={handleTabChange}
        />
        
        {/* Main Content */}
        <div className="flex-1">
          {renderContent()}
        </div>
      </div>
    </div>
  );
};

export default Profile;
