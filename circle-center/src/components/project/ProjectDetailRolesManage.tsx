import { projectApi } from '@/api/manager/project';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { UserPlus } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';

interface ProjectDetailRolesManageProps {
  projectId: number;
  onRoleAdded?: () => void;
}

export default function ProjectDetailRolesManage({
  projectId,
  onRoleAdded,
}: ProjectDetailRolesManageProps) {
  const [userId, setUserId] = useState('');
  const [role, setRole] = useState<'admin' | 'editor' | 'viewer'>('viewer');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const userIdNumber = parseInt(userId.trim());
    if (!userIdNumber || isNaN(userIdNumber)) {
      toast.error('Please enter a valid user ID');
      return;
    }

    try {
      setLoading(true);
      await projectApi.assignRole(projectId, {
        target_user_id: userIdNumber,
        role,
      });
      toast.success('Role assigned successfully');
      setUserId('');
      setRole('viewer');
      onRoleAdded?.();
    } catch (error) {
      console.error('Failed to assign role:', error);
      toast.error('Failed to assign role. Please check if the user exists.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className='flex items-center gap-2 text-sm'>
          <UserPlus className='h-4 w-4' />
          Add Collaborator
        </CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className='space-y-4'>
          <div className='space-y-2'>
            <Label htmlFor='userId'>User ID</Label>
            <Input
              id='userId'
              type='number'
              placeholder='Enter user ID'
              value={userId}
              onChange={e => setUserId(e.target.value)}
              required
            />
            <p className='text-xs text-muted-foreground'>
              Enter the numeric ID of the user you want to add
            </p>
          </div>

          <div className='space-y-2'>
            <Label htmlFor='role'>Role</Label>
            <Select
              value={role}
              onValueChange={value => setRole(value as typeof role)}
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value='admin'>Admin</SelectItem>
                <SelectItem value='editor'>Editor</SelectItem>
                <SelectItem value='viewer'>Viewer</SelectItem>
              </SelectContent>
            </Select>
            <p className='text-xs text-muted-foreground'>
              Select the role for this collaborator
            </p>
          </div>

          <Button type='submit' disabled={loading} className='w-full'>
            {loading ? 'Adding...' : 'Add Collaborator'}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
