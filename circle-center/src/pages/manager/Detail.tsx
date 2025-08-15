import ProjectDetailIconPanel from '@/components/project/ProjectDetailIconPanel';
import ProjectDetailSidebarRight from '@/components/project/ProjectDetailSidebarRight';
import ProjectDetailTokenPanel from '@/components/project/ProjectDetailTokenPanel';
import { Separator } from '@/components/ui/separator';
import { PanelRightCloseIcon, PanelRightIcon } from 'lucide-react';
import { useState } from 'react';
import { useParams } from 'react-router-dom';

export default function ManagerProjectDetail() {
  const params = useParams();
  const id = Number(params.id || 0);
  const [isSidebarOpen, setIsSidebarOpen] = useState(true);

  if (!id) return <div className='p-4'>Invalid project id</div>;

  const toggleSidebar = () => {
    setIsSidebarOpen(!isSidebarOpen);
  };

  return (
    <div className='flex flex-col lg:flex-row h-full'>
      <div className='flex-1 overflow-auto min-w-0'>
        <div className='w-full p-4 lg:p-6 space-y-6'>
          <div className='flex items-center justify-between'>
            <div className='hidden lg:block'>
              <button
                onClick={toggleSidebar}
                className='p-2 rounded-md hover:bg-accent hover:text-accent-foreground transition-colors'
                title={isSidebarOpen ? 'Hide sidebar' : 'Show sidebar'}
              >
                {isSidebarOpen ? (
                  <PanelRightCloseIcon className='h-4 w-4' />
                ) : (
                  <PanelRightIcon className='h-4 w-4' />
                )}
              </button>
            </div>
          </div>

          <ProjectDetailIconPanel projectId={id} />

          <div className='lg:hidden'>
            <Separator />
            <div>
              <ProjectDetailTokenPanel projectId={id} />
            </div>
          </div>
        </div>
      </div>

      {/* desktop */}
      <div
        className={`hidden lg:block transition-all duration-300 ease-in-out ${
          isSidebarOpen ? 'w-80' : 'w-0 overflow-hidden'
        }`}
      >
        <ProjectDetailSidebarRight projectId={id} />
      </div>
    </div>
  );
}
