import { getIconUrl, revokeIconUrl } from '@/api/manager/iconio';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { LoadingSpinner } from '@/components/ui/loading-spinner';
import { AlertCircle, Upload } from 'lucide-react';
import { useEffect, useState } from 'react';

interface ProjectDetailIconDrawableProps {
  projectId: number;
  drawable: string;
  format?: string;
  className?: string;
  size?: number;
  showUploadHint?: boolean;
  onIconStatusChange?: (hasFile: boolean) => void;
}

export default function ProjectDetailIconDrawable({
  projectId,
  drawable,
  format = 'png',
  className = '',
  size = 64,
  showUploadHint = false,
  onIconStatusChange,
}: ProjectDetailIconDrawableProps) {
  const [imageUrl, setImageUrl] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isFileNotFound, setIsFileNotFound] = useState(false);

  // Construct the relative path based on project ID and drawable
  const relPath = `icons/${projectId}/${drawable}.${format}`;

  useEffect(() => {
    let isMounted = true;

    const loadIcon = async () => {
      if (!projectId || !drawable) {
        setError('Missing project ID or drawable');
        return;
      }

      setLoading(true);
      setError(null);
      setIsFileNotFound(false);

      try {
        const url = await getIconUrl(relPath);
        if (isMounted) {
          setImageUrl(url);
          onIconStatusChange?.(true); // File exists
        }
      } catch (err) {
        if (isMounted) {
          const errorMessage =
            err instanceof Error ? err.message : 'Failed to load icon';
          setError(errorMessage);

          if (errorMessage === 'ICON_FILE_NOT_FOUND') {
            setIsFileNotFound(true);
            onIconStatusChange?.(false);
          } else {
            onIconStatusChange?.(false);
          }
        }
      } finally {
        if (isMounted) {
          setLoading(false);
        }
      }
    };

    loadIcon();

    return () => {
      isMounted = false;
      // Clean up blob URL when component unmounts
      if (imageUrl) {
        revokeIconUrl(imageUrl);
      }
    };
  }, [projectId, drawable, format, relPath, onIconStatusChange]);

  if (loading) {
    return (
      <div
        className={`flex items-center justify-center bg-muted rounded ${className}`}
        style={{ width: size, height: size }}
      >
        <LoadingSpinner size={size / 3} />
      </div>
    );
  }

  if (isFileNotFound) {
    return (
      <div
        className={`flex flex-col items-center justify-center bg-muted rounded text-muted-foreground text-xs ${className}`}
        style={{ width: size, height: size }}
      >
        <Upload className='h-6 w-6 mb-1' />
        <div className='text-center'>
          <div>No Icon File</div>
          {showUploadHint && (
            <div className='text-xs opacity-75 mt-1'>Upload required</div>
          )}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div
        className={`flex items-center justify-center bg-muted rounded ${className}`}
        style={{ width: size, height: size }}
      >
        <Alert variant='destructive' className='p-2'>
          <AlertCircle className='h-4 w-4' />
          <AlertDescription className='text-xs'>{error}</AlertDescription>
        </Alert>
      </div>
    );
  }

  if (!imageUrl) {
    return (
      <div
        className={`flex items-center justify-center bg-muted rounded text-muted-foreground text-xs ${className}`}
        style={{ width: size, height: size }}
      >
        No Icon
      </div>
    );
  }

  return (
    <img
      src={imageUrl}
      alt={`Icon: ${drawable}`}
      className={`object-contain rounded ${className}`}
      style={{ width: size, height: size }}
      onError={() => setError('Failed to load image')}
    />
  );
}
