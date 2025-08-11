import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { User, LogOut } from "lucide-react";
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
  onLogout,
  className,
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const navigate = useNavigate();

  return (
    <DropdownMenu open={isOpen} onOpenChange={setIsOpen}>
      <DropdownMenuTrigger asChild>
        <button
          className={`rounded-full focus:outline-none transition-transform duration-150 active:scale-95 ${className}`}
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
          onClick={() => {
            navigate("/profile");
            setIsOpen(false);
          }}
          className="cursor-pointer focus:bg-gray-50"
        >
          <User className="w-4 h-4 mr-2" />
          Edit Profile
        </DropdownMenuItem>
        
        <DropdownMenuSeparator />
        
        <DropdownMenuItem
          onClick={onLogout}
          className="cursor-pointer focus:bg-gray-50 text-red-600 focus:text-red-700"
        >
          <LogOut className="w-4 h-4 mr-2" />
          Sign Out
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};
