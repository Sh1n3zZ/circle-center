import { useState } from "react";
import ImportPanel from "./ProjectDetailXmlImportPanel";
import ImportConfirm from "./ProjectDetailXmlImportConfirm";
import type { IconImportComponent } from "@/api/manager/types";

export default function ProjectDetailXmlImport({ projectId }: { projectId: number }) {
  const [components, setComponents] = useState<IconImportComponent[] | null>(null);

  if (!components) {
    return <ImportPanel onParsed={setComponents} />;
  }
  return <ImportConfirm projectId={projectId} initial={components} />;
}


