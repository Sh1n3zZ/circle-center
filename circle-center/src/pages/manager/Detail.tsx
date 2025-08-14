import { useParams } from "react-router-dom";
import ProjectDetailTokenPanel from "@/components/project/ProjectDetailTokenPanel";
import ProjectDetailXmlImport from "@/components/project/ProjectDetailXmlImport";
import { Separator } from "@/components/ui/separator";

export default function ManagerProjectDetail() {
  const params = useParams();
  const id = Number(params.id || 0);

  if (!id) return <div className="p-4">Invalid project id</div>;

  return (
    <div className="max-w-4xl mx-auto p-4 space-y-6">
      <h1 className="text-xl font-semibold">Project Tokens</h1>
      <ProjectDetailTokenPanel projectId={id} />

      <Separator />

      <h2 className="text-lg font-semibold">XML Import</h2>
      <ProjectDetailXmlImport projectId={id} />
    </div>
  );
}


