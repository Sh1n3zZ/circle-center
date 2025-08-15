import { Button } from '@/components/ui/button';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import { Eye, Upload, X } from 'lucide-react';
import { useState } from 'react';
import ProjectDetailIconDrawable from './ProjectDetailIconDrawable';

interface ProjectDetailIconDrawableViewProps {
  projectId: number;
  drawable: string;
  componentInfo: string;
  size?: number;
}

export default function ProjectDetailIconDrawableView({
  projectId,
  drawable,
  componentInfo,
}: ProjectDetailIconDrawableViewProps) {
  const [open, setOpen] = useState(false);
  const [hasIconFile, setHasIconFile] = useState<boolean | null>(null);

  const handleIconLoad = (hasFile: boolean) => {
    setHasIconFile(hasFile);
  };

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant='ghost'
          size='sm'
          className='h-8 w-8 p-0 hover:bg-muted'
          title={`View icon: ${drawable}`}
        >
          <Eye className='h-4 w-4' />
          <span className='sr-only'>View icon</span>
        </Button>
      </PopoverTrigger>
      <PopoverContent className='w-96 p-4' align='start' sideOffset={8}>
        <div className='flex items-center justify-between mb-4'>
          <h4 className='font-medium leading-none'>Icon Preview</h4>
          <Button
            variant='ghost'
            size='sm'
            className='h-6 w-6 p-0'
            onClick={() => setOpen(false)}
          >
            <X className='h-4 w-4' />
            <span className='sr-only'>Close</span>
          </Button>
        </div>
        <div className='flex flex-col items-center space-y-4'>
          <div className='text-sm text-muted-foreground text-center w-full'>
            <div className='break-all'>
              <strong>Drawable:</strong> {drawable}
            </div>
            <div className='break-all'>
              <strong>Component:</strong> {componentInfo}
            </div>
          </div>
          <ProjectDetailIconDrawable
            projectId={projectId}
            drawable={drawable}
            size={128}
            className='border'
            showUploadHint={true}
            onIconStatusChange={handleIconLoad}
          />
          {hasIconFile === false && (
            <div className='text-xs text-muted-foreground text-center w-full'>
              <div className='flex items-center justify-center gap-1 mb-1'>
                <Upload className='h-3 w-3 flex-shrink-0' />
                <span>Icon file not uploaded yet</span>
              </div>
              <div className='break-words'>
                To display this icon, upload an image file for the drawable "
                {drawable}"
              </div>
            </div>
          )}
        </div>
      </PopoverContent>
    </Popover>
  );
}
