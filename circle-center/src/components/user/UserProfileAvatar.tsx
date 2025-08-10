import React from "react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { avatarApi } from "@/api/user/avatar";

interface UserProfileAvatarProps {
  /**
   * User's display name or username for fallback
   */
  displayName?: string;
  /**
   * Avatar path from backend (e.g., "avatars/2025/01/uuid_xxx.jpg")
   */
  avatarPath?: string;
  /**
   * Size of the avatar in pixels
   */
  size?: number;
  /**
   * Image quality (1-100)
   */
  quality?: number;
  /**
   * Additional CSS classes
   */
  className?: string;
}

export const UserProfileAvatar: React.FC<UserProfileAvatarProps> = ({
  displayName,
  avatarPath,
  size = 40,
  quality = 85,
  className,
}) => {
  // Generate fallback initials from display name
  const getFallbackText = (name?: string): string => {
    if (!name) return "U";
    
    const words = name.trim().split(/\s+/);
    if (words.length === 1) {
      return words[0].charAt(0).toUpperCase();
    }
    
    return (words[0].charAt(0) + words[words.length - 1].charAt(0)).toUpperCase();
  };

  // Get avatar URL if path is provided
  const avatarUrl = avatarPath ? avatarApi.getAvatarUrl(avatarPath, size, quality) : undefined;

  return (
    <Avatar className={className} style={{ width: size, height: size }}>
      {avatarUrl && (
        <AvatarImage 
          src={avatarUrl} 
          alt={`${displayName || "User"} avatar`}
          style={{ width: size, height: size }}
        />
      )}
      <AvatarFallback 
        style={{ 
          width: size, 
          height: size,
          fontSize: Math.max(12, size * 0.4),
          userSelect: "none",
          WebkitUserSelect: "none",
          MozUserSelect: "none",
          msUserSelect: "none",
        }}
      >
        {getFallbackText(displayName)}
      </AvatarFallback>
    </Avatar>
  );
};
