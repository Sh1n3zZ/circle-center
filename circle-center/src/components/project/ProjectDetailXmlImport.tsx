import type { IconImportComponent } from '@/api/manager/types';
import { useState } from 'react';
import ImportConfirm from './ProjectDetailXmlImportConfirm';
import ImportPanel from './ProjectDetailXmlImportPanel';

export default function ProjectDetailXmlImport({
  projectId,
}: {
  projectId: number;
}) {
  const [components, setComponents] = useState<IconImportComponent[] | null>(
    null
  );

  if (!components) {
    return <ImportPanel onParsed={setComponents} />;
  }
  return (
    <ImportConfirm
      projectId={projectId}
      initial={components}
      onClose={() => setComponents(null)}
    />
  );
}
