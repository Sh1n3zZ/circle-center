import {
  SidebarGroup,
  SidebarGroupAction,
  SidebarGroupContent,
  SidebarGroupLabel,
} from '@/components/ui/sidebar';
import { ChevronDown, ChevronRight, FileText, Key } from 'lucide-react';
import { useState } from 'react';
import ProjectDetailTokenPanel from './ProjectDetailTokenPanel';
import ProjectDetailXmlImport from './ProjectDetailXmlImport';

interface ProjectDetailSidebarRightProps {
  projectId: number;
}

export default function ProjectDetailSidebarRight({
  projectId,
}: ProjectDetailSidebarRightProps) {
  const [isTokensExpanded, setIsTokensExpanded] = useState(true);
  const [isXmlExpanded, setIsXmlExpanded] = useState(false);

  const toggleTokensExpanded = () => {
    setIsTokensExpanded(!isTokensExpanded);
  };

  const toggleXmlExpanded = () => {
    setIsXmlExpanded(!isXmlExpanded);
  };

  return (
    <div className='sticky top-0 border-l w-80 min-w-80 h-full bg-background'>
      <div className='border-b h-16 flex items-center px-4'>
        <h2 className='text-lg font-semibold'>Project Tools</h2>
      </div>
      <div className='p-2 overflow-auto h-[calc(100%-4rem)] space-y-2'>
        <SidebarGroup>
          <SidebarGroupLabel
            className='flex items-center justify-between cursor-pointer'
            onClick={toggleTokensExpanded}
          >
            <div className='flex items-center gap-2'>
              <Key className='h-4 w-4' />
              <span>Project Tokens</span>
            </div>
            <SidebarGroupAction
              onClick={e => {
                e.stopPropagation();
                toggleTokensExpanded();
              }}
            >
              {isTokensExpanded ? (
                <ChevronDown className='h-4 w-4' />
              ) : (
                <ChevronRight className='h-4 w-4' />
              )}
            </SidebarGroupAction>
          </SidebarGroupLabel>
          <SidebarGroupContent
            className={`transition-all duration-300 ease-in-out ${
              isTokensExpanded
                ? 'max-h-[800px] opacity-100 mt-2'
                : 'max-h-0 opacity-0 overflow-hidden'
            }`}
          >
            <ProjectDetailTokenPanel projectId={projectId} />
          </SidebarGroupContent>
        </SidebarGroup>
        <SidebarGroup>
          <SidebarGroupLabel
            className='flex items-center justify-between cursor-pointer'
            onClick={toggleXmlExpanded}
          >
            <div className='flex items-center gap-2'>
              <FileText className='h-4 w-4' />
              <span>XML Import</span>
            </div>
            <SidebarGroupAction
              onClick={e => {
                e.stopPropagation();
                toggleXmlExpanded();
              }}
            >
              {isXmlExpanded ? (
                <ChevronDown className='h-4 w-4' />
              ) : (
                <ChevronRight className='h-4 w-4' />
              )}
            </SidebarGroupAction>
          </SidebarGroupLabel>
          <SidebarGroupContent
            className={`transition-all duration-300 ease-in-out ${
              isXmlExpanded
                ? 'max-h-[800px] opacity-100 mt-2'
                : 'max-h-0 opacity-0 overflow-hidden'
            }`}
          >
            <ProjectDetailXmlImport projectId={projectId} />
          </SidebarGroupContent>
        </SidebarGroup>
      </div>
    </div>
  );
}
