import React, { useState, useRef } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Pencil } from "lucide-react";
import { avatarApi } from "@/api/user/avatar";
import { LoadingSpinner } from "@/components/ui/loading-spinner";
import { toast } from "sonner";

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
   * Whether avatar editing is enabled
   */
  editable?: boolean;
  /**
   * Callback when avatar is successfully uploaded
   */
  onAvatarChange?: (newAvatarPath: string) => void;
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
  editable = false,
  onAvatarChange,
  className,
}) => {
  const [isUploading, setIsUploading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const getFallbackText = (name?: string): string => {
    if (!name) return "U";
    
    const words = name.trim().split(/\s+/);
    if (words.length === 1) {
      return words[0].charAt(0).toUpperCase();
    }
    
    return (words[0].charAt(0) + words[words.length - 1].charAt(0)).toUpperCase();
  };

  // get avatar URL if path is provided
  const avatarUrl = avatarPath ? avatarApi.getAvatarUrl(avatarPath, size, quality) : undefined;

  const handleFileSelect = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    if (!file.type.startsWith('image/')) {
      toast.error("Please select a valid image file");
      return;
    }

    if (file.size > 5 * 1024 * 1024) {
      toast.error("Image file size must be less than 5MB");
      return;
    }

    setIsUploading(true);
    try {
      const result = await avatarApi.uploadAvatar(file);
      onAvatarChange?.(result.path);
      toast.success("Avatar updated successfully");
    } catch (error) {
      console.error("Failed to upload avatar:", error);
      toast.error("Failed to upload avatar. Please try again.");
    } finally {
      setIsUploading(false);
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
    }
  };

  const handleUploadClick = () => {
    if (editable && !isUploading) {
      fileInputRef.current?.click();
    }
  };

  return (
    <div className={`relative inline-block ${className}`}>
      <Avatar 
        className="relative"
        style={{ width: size, height: size }}
      >
        {avatarUrl && (
          <AvatarImage 
            src={avatarUrl} 
            alt={`${displayName || "User"} avatar`}
            style={{ width: size, height: size }}
            onError={(e) => {
              console.error('Avatar image failed to load:', avatarUrl, e);
            }}
            onLoad={() => {
              console.log('Avatar image loaded successfully:', avatarUrl);
            }}
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

      {editable && (
        <div className="absolute inset-0">
            <Button
              variant="secondary"
              size="sm"
              onClick={handleUploadClick}
              disabled={isUploading}
              className="absolute inset-0 w-full h-full rounded-full p-0 bg-black/20 hover:bg-black/30 transition-all duration-200 opacity-0 hover:opacity-100 group cursor-pointer"
            >
            {isUploading ? (
              <LoadingSpinner size={24} className="text-white" />
                         ) : (
               <div className="flex items-center justify-center text-white">
                 <Pencil className="w-5 h-5" />
               </div>
             )}
          </Button>
        </div>
      )}

      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        onChange={handleFileSelect}
        className="hidden"
        disabled={isUploading}
      />
    </div>
  );
};
