import type { IconModel } from '@/api/manager/types';
import { Button } from '@/components/ui/button';
import { Edit, Trash2 } from 'lucide-react';

interface ProjectDetailIconTableActionsProps {
  icon: IconModel;
  onEdit?: (icon: IconModel) => void;
  onDelete?: (icon: IconModel) => void;
}

export default function ProjectDetailIconTableActions({
  icon,
  onEdit,
  onDelete,
}: ProjectDetailIconTableActionsProps) {
  return (
    <div className='flex items-center space-x-1'>
      {onEdit && (
        <Button
          variant='ghost'
          size='sm'
          className='h-8 w-8 p-0'
          onClick={() => onEdit(icon)}
          title='Edit icon'
        >
          <Edit className='h-4 w-4' />
          <span className='sr-only'>Edit icon</span>
        </Button>
      )}
      {onDelete && (
        <Button
          variant='ghost'
          size='sm'
          className='h-8 w-8 p-0 text-red-600 hover:text-red-700 hover:bg-red-50'
          onClick={() => onDelete(icon)}
          title='Delete icon'
        >
          <Trash2 className='h-4 w-4' />
          <span className='sr-only'>Delete icon</span>
        </Button>
      )}
    </div>
  );
}
