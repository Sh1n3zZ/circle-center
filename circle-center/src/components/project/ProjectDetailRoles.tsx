import { useState } from 'react';
import ProjectDetailRolesIDs from './ProjectDetailRolesIDs';
import ProjectDetailRolesManage from './ProjectDetailRolesManage';

interface ProjectDetailRolesProps {
  projectId: number;
}

export default function ProjectDetailRoles({
  projectId,
}: ProjectDetailRolesProps) {
  const [refreshKey, setRefreshKey] = useState(0);

  const handleRoleUpdate = () => {
    // Trigger refresh of the roles list
    setRefreshKey(prev => prev + 1);
  };

  return (
    <div className='space-y-4'>
      <ProjectDetailRolesManage
        projectId={projectId}
        onRoleAdded={handleRoleUpdate}
      />
      <div key={refreshKey}>
        <ProjectDetailRolesIDs
          projectId={projectId}
          onRoleUpdate={handleRoleUpdate}
        />
      </div>
    </div>
  );
}
