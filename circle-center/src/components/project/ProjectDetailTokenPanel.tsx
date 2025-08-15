import ProjectDetailTokenCreate from '@/components/project/ProjectDetailTokenCreate';
import ProjectDetailTokenList from '@/components/project/ProjectDetailTokenList';
import { useCallback, useState } from 'react';

export default function ProjectDetailTokenPanel({
  projectId,
}: {
  projectId: number;
}) {
  const [refresh, setRefresh] = useState(0);

  const handleCreated = useCallback(() => {
    setRefresh(v => v + 1);
  }, []);

  return (
    <div className='space-y-4'>
      <ProjectDetailTokenCreate
        projectId={projectId}
        onCreated={handleCreated}
      />
      <ProjectDetailTokenList projectId={projectId} refreshSignal={refresh} />
    </div>
  );
}
