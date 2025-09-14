import { projectApi } from '@/api/manager/project';
import type { CollaboratorInfo } from '@/api/manager/types';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Skeleton } from '@/components/ui/skeleton';
import { UserProfileAvatar } from '@/components/user/UserProfileAvatar';
import { Eye, MoreHorizontal, Shield, Trash2, UserCheck } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';

interface ProjectDetailRolesIDsProps {
  projectId: number;
  onRoleUpdate?: () => void;
}

const getRoleIcon = (role: string) => {
  switch (role) {
    case 'owner':
      return <Shield className='h-3 w-3' />;
    case 'admin':
      return <UserCheck className='h-3 w-3' />;
    case 'editor':
      return <UserCheck className='h-3 w-3' />;
    case 'viewer':
      return <Eye className='h-3 w-3' />;
    default:
      return <UserCheck className='h-3 w-3' />;
  }
};

const getRoleBadgeVariant = (role: string) => {
  switch (role) {
    case 'owner':
      return 'default';
    case 'admin':
      return 'secondary';
    case 'editor':
      return 'outline';
    case 'viewer':
      return 'outline';
    default:
      return 'outline';
  }
};

export default function ProjectDetailRolesIDs({
  projectId,
  onRoleUpdate,
}: ProjectDetailRolesIDsProps) {
  const [collaborators, setCollaborators] = useState<CollaboratorInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [updating, setUpdating] = useState<number | null>(null);

  const fetchCollaborators = async () => {
    try {
      setLoading(true);
      const response = await projectApi.listProjectRoles(projectId);
      setCollaborators(response.data);
    } catch (error) {
      console.error('Failed to fetch collaborators:', error);
      toast.error('Failed to load collaborators');
    } finally {
      setLoading(false);
    }
  };

  const handleRoleChange = async (userId: number, newRole: string) => {
    try {
      setUpdating(userId);
      await projectApi.assignRole(projectId, {
        target_user_id: userId,
        role: newRole as 'admin' | 'editor' | 'viewer',
      });
      toast.success('Role updated successfully');
      await fetchCollaborators();
      onRoleUpdate?.();
    } catch (error) {
      console.error('Failed to update role:', error);
      toast.error('Failed to update role');
    } finally {
      setUpdating(null);
    }
  };

  const handleRemoveCollaborator = async (userId: number) => {
    try {
      setUpdating(userId);
      await projectApi.removeCollaborator(projectId, userId);
      toast.success('Collaborator removed successfully');
      await fetchCollaborators();
      onRoleUpdate?.();
    } catch (error) {
      console.error('Failed to remove collaborator:', error);
      toast.error('Failed to remove collaborator');
    } finally {
      setUpdating(null);
    }
  };

  useEffect(() => {
    fetchCollaborators();
  }, [projectId]);

  if (loading) {
    return (
      <div className='space-y-2'>
        {[...Array(3)].map((_, i) => (
          <Card key={i}>
            <CardContent className='p-3'>
              <div className='flex items-center justify-between'>
                <div className='flex items-center gap-2'>
                  <Skeleton className='h-8 w-8 rounded-full' />
                  <div>
                    <Skeleton className='h-4 w-20 mb-1' />
                    <Skeleton className='h-3 w-16' />
                  </div>
                </div>
                <Skeleton className='h-6 w-16' />
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  return (
    <div className='space-y-2'>
      {collaborators.map(collaborator => (
        <Card key={collaborator.user_id}>
          <CardContent className='p-3'>
            <div className='flex items-center justify-between'>
              <div className='flex items-center gap-2'>
                <UserProfileAvatar
                  displayName={
                    collaborator.display_name || collaborator.username
                  }
                  avatarPath={collaborator.avatar_url}
                  size={32}
                />
                <div>
                  <div className='font-medium text-sm'>
                    {collaborator.display_name ||
                      collaborator.username ||
                      `User ${collaborator.user_id}`}
                  </div>
                  {collaborator.username && (
                    <div className='text-xs text-muted-foreground'>
                      @{collaborator.username}
                    </div>
                  )}
                </div>
              </div>
              <div className='flex items-center gap-2'>
                <Badge
                  variant={getRoleBadgeVariant(collaborator.role)}
                  className='flex items-center gap-1'
                >
                  {getRoleIcon(collaborator.role)}
                  {collaborator.role}
                </Badge>
                {collaborator.role !== 'owner' && (
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button
                        variant='ghost'
                        size='sm'
                        className='h-8 w-8 p-0'
                        disabled={updating === collaborator.user_id}
                      >
                        <MoreHorizontal className='h-4 w-4' />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align='end'>
                      {['admin', 'editor', 'viewer'].map(role => (
                        <DropdownMenuItem
                          key={role}
                          onClick={() =>
                            handleRoleChange(collaborator.user_id, role)
                          }
                          disabled={role === collaborator.role}
                          className='flex items-center gap-2'
                        >
                          {getRoleIcon(role)}
                          Change to {role}
                        </DropdownMenuItem>
                      ))}
                      <DropdownMenuItem
                        onClick={() =>
                          handleRemoveCollaborator(collaborator.user_id)
                        }
                        className='flex items-center gap-2 text-destructive focus:text-destructive'
                      >
                        <Trash2 className='h-3 w-3' />
                        Remove from project
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                )}
              </div>
            </div>
          </CardContent>
        </Card>
      ))}
      {collaborators.length === 0 && (
        <div className='text-center py-8 text-muted-foreground'>
          No collaborators found
        </div>
      )}
    </div>
  );
}
