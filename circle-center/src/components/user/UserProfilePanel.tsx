import React, { useState, useEffect } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { UserProfileAvatar } from "./UserProfileAvatar";
import { profileApi } from "@/api/user/profile";
import { toast } from "sonner";
import { LoadingSpinner } from "@/components/ui/loading-spinner";
import type { GetUserProfileResponse, UpdateUserProfileRequest } from "@/api/user/types";
import jstz from "jstz";

interface UserProfilePanelProps {
  /**
   * Additional CSS classes
   */
  className?: string;
}

export const UserProfilePanel: React.FC<UserProfilePanelProps> = ({ className }) => {
  const [profile, setProfile] = useState<GetUserProfileResponse['data'] | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [detectedLocale, setDetectedLocale] = useState<string>("");
  const [detectedTimezone, setDetectedTimezone] = useState<string>("");
  const [formData, setFormData] = useState<UpdateUserProfileRequest>({
    display_name: "",
    phone: "",
  });

  useEffect(() => {
    loadUserProfile();
    detectUserLocaleAndTimezone();
  }, []);

  const detectUserLocaleAndTimezone = () => {
    const detectedLocale = navigator.language || navigator.languages?.[0] || "en";
    setDetectedLocale(detectedLocale);

    const detectedTimezone = jstz.determine().name() || "UTC";
    setDetectedTimezone(detectedTimezone);
  };

  const loadUserProfile = async () => {
    try {
      setIsLoading(true);
      const response = await profileApi.getUserProfile();
      setProfile(response.data);
      setFormData({
        display_name: response.data.display_name || "",
        phone: response.data.phone || "",
      });
    } catch (error) {
      console.error("Failed to load user profile:", error);
      toast.error("Failed to load user profile");
    } finally {
      setIsLoading(false);
    }
  };

  const handleInputChange = (field: keyof UpdateUserProfileRequest, value: string) => {
    setFormData(prev => ({
      ...prev,
      [field]: value,
    }));
  };

  const handleSave = async () => {
    try {
      setIsSaving(true);
      const payload: UpdateUserProfileRequest = {
        display_name: formData.display_name,
        phone: formData.phone,
      };

      if (!profile || profile.locale !== detectedLocale) {
        payload.locale = detectedLocale;
      }
      if (!profile || profile.timezone !== detectedTimezone) {
        payload.timezone = detectedTimezone;
      }

      const response = await profileApi.updateUserProfile(payload);

      if (profile) {
        setProfile({
          ...response.data,
          created_at: profile.created_at,
        });
      } else {
        setProfile({
          ...response.data,
          created_at: new Date().toISOString(),
        } as GetUserProfileResponse['data']);
      }

      toast.success("Profile updated successfully");
    } catch (error) {
      console.error("Failed to update profile:", error);
      toast.error("Failed to update profile");
    } finally {
      setIsSaving(false);
    }
  };

  const handleAvatarChange = (newAvatarPath: string) => {
    if (profile) {
      setProfile({
        ...profile,
        avatar_url: newAvatarPath,
      });
    }
  };

  if (isLoading) {
    return (
      <Card className={className}>
        <CardHeader>
          <CardTitle>Profile Information</CardTitle>
          <CardDescription>Manage your account settings and preferences.</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-center py-8">
            <LoadingSpinner size={32} className="text-gray-900" />
          </div>
        </CardContent>
      </Card>
    );
  }

  if (!profile) {
    return (
      <Card className={className}>
        <CardContent>
          <div className="text-center py-8">
            <p className="text-gray-500">Failed to load profile information</p>
            <Button onClick={loadUserProfile} className="mt-4">
              Retry
            </Button>
          </div>
        </CardContent>
      </Card>
    );
  }

  const renderChangedValue = (storageValue: string, detectValue: string) => {
    const isSame = storageValue === detectValue;
    if (isSame) {
      return (
        <Input
          value={storageValue}
          disabled
          className="bg-gray-50"
        />
      );
    }
    return (
      <div className="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-50 text-gray-900 min-h-[40px] flex items-center">
        <del className="text-gray-400">{storageValue}</del>
        <span className="ml-2">{detectValue}</span>
      </div>
    );
  };

  return (
    <Card className={className}>
      <CardHeader>
        <CardTitle>Profile Information</CardTitle>
        <CardDescription>Manage your account settings and preferences.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Avatar Section */}
        <div className="flex items-center space-x-4">
          <UserProfileAvatar
            displayName={profile.display_name || profile.username}
            avatarPath={profile.avatar_url}
            size={80}
            editable={true}
            onAvatarChange={handleAvatarChange}
          />
          <div>
            <h3 className="text-lg font-medium">Profile Picture</h3>
            <p className="text-sm text-gray-500">
              Click on the avatar to upload a new image. JPG, PNG or GIF up to 5MB.
            </p>
          </div>
        </div>

        {/* Profile Form */}
        <div className="grid gap-4">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="username">Username</Label>
              <Input
                id="username"
                value={profile.username}
                disabled
                className="bg-gray-50"
              />
              <p className="text-xs text-gray-500">Username cannot be changed</p>
            </div>
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                value={profile.email}
                disabled
                className="bg-gray-50"
              />
              <p className="text-xs text-gray-500">Email cannot be changed</p>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="display_name">Display Name</Label>
            <Input
              id="display_name"
              value={formData.display_name}
              onChange={(e) => handleInputChange("display_name", e.target.value)}
              placeholder="Enter your display name"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="phone">Phone Number</Label>
            <Input
              id="phone"
              value={formData.phone}
              onChange={(e) => handleInputChange("phone", e.target.value)}
              placeholder="Enter your phone number"
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="locale">Language</Label>
              {renderChangedValue(profile.locale, detectedLocale)}
              <p className="text-xs text-gray-500">Detected from browser settings</p>
            </div>
            <div className="space-y-2">
              <Label htmlFor="timezone">Timezone</Label>
              {renderChangedValue(profile.timezone, detectedTimezone)}
              <p className="text-xs text-gray-500">Detected from system settings</p>
            </div>
          </div>

          <div className="pt-4">
            <Button 
              onClick={handleSave} 
              disabled={isSaving}
              className="w-full sm:w-auto"
            >
              {isSaving ? "Saving..." : "Save Changes"}
            </Button>
          </div>
        </div>

        {/* Account Info */}
        <div className="border-t pt-6">
          <h4 className="text-sm font-medium text-gray-900 mb-2">Account Information</h4>
          <div className="text-sm text-gray-500 space-y-1">
            <p>Member since: {new Date(profile.created_at).toLocaleDateString()}</p>
            <p>Last updated: {new Date(profile.updated_at).toLocaleDateString()}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};
