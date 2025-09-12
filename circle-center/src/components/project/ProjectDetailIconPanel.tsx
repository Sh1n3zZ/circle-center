import type { IconModel } from '@/api/manager/types';
import { toast } from 'sonner';
import ProjectDetailIconCard from './ProjectDetailIconCard';

interface ProjectDetailIconPanelProps {
  projectId: number;
}

export default function ProjectDetailIconPanel({
  projectId,
}: ProjectDetailIconPanelProps) {
  const handleAddIcon = () => {
    // TODO: Implement add icon functionality
    toast.info('Add icon functionality coming soon');
  };

  const handleEditIcon = (icon: IconModel) => {
    // TODO: Implement edit icon functionality
    toast.info(`Edit icon: ${icon.name}`);
  };

  const handleDeleteIcon = async (icon: IconModel) => {
    // Confirmation is handled in the card component
    // This is just a callback for additional logic if needed
    console.log('Icon deleted:', icon.id);
  };

  return (
    <div className='space-y-4'>
      <div className='flex items-center justify-between'>
        <div>
          <h2 className='text-lg font-semibold'>Icon Management</h2>
          <p className='text-sm text-muted-foreground'>
            Manage icons for this project. View, edit, and delete existing
            icons.
          </p>
        </div>
      </div>

      <ProjectDetailIconCard
        projectId={projectId}
        onAddIcon={handleAddIcon}
        onEditIcon={handleEditIcon}
        onDeleteIcon={handleDeleteIcon}
      />
    </div>
  );
}
