import React from "react";
import { cn } from "@/lib/utils";
import { User, Shield, Settings } from "lucide-react";
import { ScrollArea } from "@/components/ui/scroll-area";

export type ProfileTab = "profile" | "security";

interface UserProfileSideBarProps {
  /**
   * Currently active tab
   */
  activeTab: ProfileTab;
  /**
   * Callback when tab changes
   */
  onTabChange: (tab: ProfileTab) => void;
  /**
   * Additional CSS classes
   */
  className?: string;
}

export const UserProfileSideBar: React.FC<UserProfileSideBarProps> = ({
  activeTab,
  onTabChange,
  className,
}) => {
  const tabs = [
    {
      id: "profile" as ProfileTab,
      label: "Profile",
      icon: <User className="w-5 h-5" />,
      description: "Manage your personal information",
    },
    {
      id: "security" as ProfileTab,
      label: "Security",
      icon: <Shield className="w-5 h-5" />,
      description: "Account security and privacy",
    },
  ];

  return (
    <div className={cn("w-64 h-screen bg-white border-r border-gray-200 shrink-0", className)}>
      <ScrollArea className="h-full">
        <div className="p-6">
        <div className="flex items-center space-x-3 mb-6">
          <Settings className="w-6 h-6 text-gray-600" />
          <h2 className="text-lg font-semibold text-gray-900">Account Settings</h2>
        </div>
        
        <nav className="space-y-1">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => onTabChange(tab.id)}
              className={cn(
                "w-full flex items-start px-4 py-3 text-sm font-medium rounded-lg transition-all duration-200 text-left",
                activeTab === tab.id
                  ? "bg-blue-50 text-blue-700 border border-blue-200 shadow-sm"
                  : "text-gray-600 hover:text-gray-900 hover:bg-gray-50 border border-transparent"
              )}
            >
              <span className={cn(
                "mr-3 mt-0.5",
                activeTab === tab.id ? "text-blue-600" : "text-gray-400"
              )}>
                {tab.icon}
              </span>
              <div className="flex-1">
                <div className="font-medium">{tab.label}</div>
                <div className={cn(
                  "text-xs mt-0.5",
                  activeTab === tab.id ? "text-blue-600" : "text-gray-400"
                )}>
                  {tab.description}
                </div>
              </div>
            </button>
          ))}
        </nav>
        </div>
      </ScrollArea>
    </div>
  );
};
