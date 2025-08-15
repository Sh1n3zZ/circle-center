import { iconApi } from '@/api/manager/icon';
import type { IconModel, UpdateIconRequest } from '@/api/manager/types';
import { Alert, AlertDescription } from '@/components/ui/alert';
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { LoadingSpinner } from '@/components/ui/loading-spinner';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { AlertCircle, Save, X } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import ProjectDetailIconDrawable from './ProjectDetailIconDrawable';
import ProjectDetailUploadDrawable from './ProjectDetailUploadDrawable';

interface ProjectDetailIconEditProps {
  projectId: number;
  icon: IconModel | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSave?: (updatedIcon: IconModel) => void;
}

export default function ProjectDetailIconEdit({
  projectId,
  icon,
  open,
  onOpenChange,
  onSave,
}: ProjectDetailIconEditProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [iconRefreshKey, setIconRefreshKey] = useState(0);
  const [formData, setFormData] = useState<UpdateIconRequest>({
    name: '',
    pkg: '',
    componentInfo: '',
    drawable: '',
    status: 'pending',
  });

  // Initialize form data when icon changes
  useEffect(() => {
    if (icon) {
      setFormData({
        name: icon.name,
        pkg: icon.pkg,
        componentInfo: icon.componentInfo,
        drawable: icon.drawable,
        status: icon.status,
      });
      setIconRefreshKey(prev => prev + 1);
    }
  }, [icon]);

  const handleInputChange = (field: keyof UpdateIconRequest, value: string) => {
    setFormData(prev => ({
      ...prev,
      [field]: value,
    }));

    // Refresh icon display when drawable changes
    if (field === 'drawable') {
      setIconRefreshKey(prev => prev + 1);
    }
  };

  const handleSave = async () => {
    if (!icon) return;

    setLoading(true);
    setError(null);

    try {
      const response = await iconApi.update(projectId, icon.id, formData);

      if (response.success) {
        toast.success('Icon updated successfully');
        onSave?.(response.data);
        onOpenChange(false);
      } else {
        throw new Error(response.message || 'Failed to update icon');
      }
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : 'Failed to update icon';
      setError(errorMessage);
      toast.error(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  const handleUploadSuccess = () => {
    toast.success('Icon file uploaded successfully');
    setIconRefreshKey(prev => prev + 1);
  };

  const handleUploadError = (error: string) => {
    toast.error(`Upload failed: ${error}`);
  };

  const handleClose = () => {
    setError(null);
    onOpenChange(false);
  };

  if (!icon) {
    return null;
  }

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent className='sm:max-w-4xl max-h-[90vh] overflow-y-auto'>
        <AlertDialogHeader>
          <AlertDialogTitle className='break-words'>
            Edit Icon: {icon.name}
          </AlertDialogTitle>
        </AlertDialogHeader>

        <div className='space-y-6'>
          {error && (
            <Alert variant='destructive'>
              <AlertCircle className='h-4 w-4' />
              <AlertDescription className='break-words'>
                {error}
              </AlertDescription>
            </Alert>
          )}

          <div className='grid grid-cols-1 lg:grid-cols-2 gap-6'>
            {/* Left column - Form fields */}
            <div className='space-y-4'>
              <div className='space-y-2'>
                <Label htmlFor='name'>Name</Label>
                <Input
                  id='name'
                  value={formData.name}
                  onChange={e => handleInputChange('name', e.target.value)}
                  placeholder='Icon name'
                  className='break-all'
                />
              </div>

              <div className='space-y-2'>
                <Label htmlFor='pkg'>Package</Label>
                <Input
                  id='pkg'
                  value={formData.pkg}
                  onChange={e => handleInputChange('pkg', e.target.value)}
                  placeholder='com.example.app'
                  className='break-all'
                />
              </div>

              <div className='space-y-2'>
                <Label htmlFor='componentInfo'>Component Info</Label>
                <Input
                  id='componentInfo'
                  value={formData.componentInfo}
                  onChange={e =>
                    handleInputChange('componentInfo', e.target.value)
                  }
                  placeholder='com.example.app/.MainActivity'
                  className='break-all'
                />
              </div>

              <div className='space-y-2'>
                <Label htmlFor='drawable'>Drawable</Label>
                <Input
                  id='drawable'
                  value={formData.drawable}
                  onChange={e => handleInputChange('drawable', e.target.value)}
                  placeholder='ic_launcher'
                  className='break-all'
                />
              </div>

              <div className='space-y-2'>
                <Label htmlFor='status'>Status</Label>
                <Select
                  value={formData.status}
                  onValueChange={value => handleInputChange('status', value)}
                >
                  <SelectTrigger>
                    <SelectValue placeholder='Select status' />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value='pending'>Pending</SelectItem>
                    <SelectItem value='in_progress'>In Progress</SelectItem>
                    <SelectItem value='published'>Published</SelectItem>
                    <SelectItem value='rejected'>Rejected</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            {/* Right column - Icon preview and upload */}
            <div className='space-y-4'>
              <div className='space-y-2'>
                <Label>Current Icon</Label>
                <div className='flex justify-center'>
                  <ProjectDetailIconDrawable
                    key={`${icon.id}-${iconRefreshKey}`}
                    projectId={projectId}
                    drawable={formData.drawable || icon.drawable}
                    size={128}
                    className='border rounded-lg'
                  />
                </div>
              </div>

              <div className='space-y-2'>
                <Label>Upload New Icon</Label>
                <ProjectDetailUploadDrawable
                  projectId={projectId}
                  componentInfo={formData.componentInfo || icon.componentInfo}
                  drawable={formData.drawable || icon.drawable}
                  onUploadSuccess={handleUploadSuccess}
                  onUploadError={handleUploadError}
                />
              </div>
            </div>
          </div>
        </div>

        <AlertDialogFooter>
          <Button variant='outline' onClick={handleClose} disabled={loading}>
            <X className='mr-2 h-4 w-4' />
            Cancel
          </Button>
          <Button onClick={handleSave} disabled={loading}>
            {loading ? (
              <>
                <LoadingSpinner size={16} />
                <span className='ml-2'>Saving...</span>
              </>
            ) : (
              <>
                <Save className='mr-2 h-4 w-4' />
                Save Changes
              </>
            )}
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
