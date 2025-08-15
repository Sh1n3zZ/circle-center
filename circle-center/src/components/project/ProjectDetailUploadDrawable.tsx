import { uploadIcon } from '@/api/manager/iconio';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { LoadingSpinner } from '@/components/ui/loading-spinner';
import { AlertCircle, CheckCircle, Upload } from 'lucide-react';
import { useRef, useState } from 'react';
import { toast } from 'sonner';

interface ProjectDetailUploadDrawableProps {
  projectId: number;
  componentInfo: string;
  drawable: string;
  onUploadSuccess?: (path: string) => void;
  onUploadError?: (error: string) => void;
}

export default function ProjectDetailUploadDrawable({
  projectId,
  componentInfo,
  drawable,
  onUploadSuccess,
  onUploadError,
}: ProjectDetailUploadDrawableProps) {
  const [uploading, setUploading] = useState(false);
  const [uploaded, setUploaded] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileSelect = async (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    const file = event.target.files?.[0];
    if (!file) return;

    // Validate file type
    const validTypes = ['image/png', 'image/jpeg', 'image/jpg', 'image/gif'];
    if (!validTypes.includes(file.type)) {
      setError('Please select a valid image file (PNG, JPEG, or GIF)');
      return;
    }

    // Validate file size (5MB limit)
    const maxSize = 5 * 1024 * 1024; // 5MB
    if (file.size > maxSize) {
      setError('File size must be less than 5MB');
      return;
    }

    setUploading(true);
    setError(null);

    try {
      const response = await uploadIcon(projectId, componentInfo, file);

      if (response.success) {
        setUploaded(true);
        toast.success('Icon uploaded successfully');
        onUploadSuccess?.(response.data.path);
      } else {
        throw new Error(response.message || 'Upload failed');
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Upload failed';
      setError(errorMessage);
      toast.error(errorMessage);
      onUploadError?.(errorMessage);
    } finally {
      setUploading(false);
    }
  };

  const handleUploadClick = () => {
    fileInputRef.current?.click();
  };

  const resetUpload = () => {
    setUploaded(false);
    setError(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  return (
    <div className='space-y-4'>
      <div className='space-y-2'>
        <Label htmlFor='icon-upload'>Upload Icon for: {drawable}</Label>
        <div className='text-sm text-muted-foreground'>
          Component: {componentInfo}
        </div>
      </div>

      <div className='flex items-center space-x-2'>
        <Button
          onClick={handleUploadClick}
          disabled={uploading}
          variant='outline'
          className='flex items-center space-x-2'
        >
          {uploading ? (
            <>
              <LoadingSpinner size={16} />
              <span>Uploading...</span>
            </>
          ) : uploaded ? (
            <>
              <CheckCircle className='h-4 w-4 text-green-600' />
              <span>Uploaded</span>
            </>
          ) : (
            <>
              <Upload className='h-4 w-4' />
              <span>Choose File</span>
            </>
          )}
        </Button>

        {uploaded && (
          <Button onClick={resetUpload} variant='ghost' size='sm'>
            Upload Another
          </Button>
        )}
      </div>

      <input
        ref={fileInputRef}
        type='file'
        accept='image/png,image/jpeg,image/jpg,image/gif'
        onChange={handleFileSelect}
        className='hidden'
        id='icon-upload'
      />

      {error && (
        <Alert variant='destructive'>
          <AlertCircle className='h-4 w-4' />
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {uploaded && (
        <Alert>
          <CheckCircle className='h-4 w-4' />
          <AlertDescription>
            Icon uploaded successfully! The file will be saved as "{drawable}
            .png"
          </AlertDescription>
        </Alert>
      )}
    </div>
  );
}
