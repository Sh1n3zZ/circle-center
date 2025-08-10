import React, { useState } from "react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { UserProfileAvatar } from "./UserProfileAvatar";

interface UserProfileNavBubbleProps {
  /**
   * User's display name
   */
  displayName?: string;
  /**
   * Avatar path from backend
   */
  avatarPath?: string;
  /**
   * Size of the avatar in pixels
   */
  size?: number;
  /**
   * Callback for edit profile action
   */
  onEditProfile?: () => void;
  /**
   * Callback for logout action
   */
  onLogout?: () => void;
  /**
   * Additional CSS classes
   */
  className?: string;
}

export const UserProfileNavBubble: React.FC<UserProfileNavBubbleProps> = ({
  displayName,
  avatarPath,
  size = 32,
  onEditProfile,
  onLogout,
  className,
}) => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <DropdownMenu open={isOpen} onOpenChange={setIsOpen}>
      <DropdownMenuTrigger asChild>
        <button
          className={`rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-all duration-200 hover:scale-105 ${className}`}
          aria-label={`${displayName || "User"} profile menu`}
        >
          <UserProfileAvatar
            displayName={displayName}
            avatarPath={avatarPath}
            size={size}
            className="cursor-pointer"
          />
        </button>
      </DropdownMenuTrigger>
      
      <DropdownMenuContent 
        align="end" 
        className="w-56 mt-2"
        sideOffset={8}
      >
        <div className="px-3 py-2 border-b border-gray-100">
          <p className="text-sm font-medium text-gray-900 truncate">
            {displayName || "User"}
          </p>
          <p className="text-xs text-gray-500 truncate">
            Profile Settings
          </p>
        </div>
        
        <DropdownMenuItem
          onClick={onEditProfile}
          className="cursor-pointer focus:bg-gray-50"
        >
          <svg
            className="w-4 h-4 mr-2"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
            />
          </svg>
          Edit Profile
        </DropdownMenuItem>
        
        <DropdownMenuSeparator />
        
        <DropdownMenuItem
          onClick={onLogout}
          className="cursor-pointer focus:bg-gray-50 text-red-600 focus:text-red-700"
        >
          <svg
            className="w-4 h-4 mr-2"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
            />
          </svg>
          Sign Out
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};
