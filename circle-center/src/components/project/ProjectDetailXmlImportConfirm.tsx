import { useState } from "react";
import { Button } from "@/components/ui/button";
import ConfirmTable from "./ProjectDetailXmlImportConfirmTable";
import type { IconImportComponent, ConfirmImportRequest } from "@/api/manager/types";
import { xmlioApi } from "@/api/manager/xmlio";
import { toast } from "sonner";

export default function ProjectDetailXmlImportConfirm({ projectId, initial }: { projectId: number; initial: IconImportComponent[] }) {
  const [data, setData] = useState<IconImportComponent[]>(initial);
  const [loading, setLoading] = useState(false);

  const handleConfirm = async () => {
    try {
      setLoading(true);
      const body: ConfirmImportRequest = { projectId, components: data };
      const res = await xmlioApi.confirmImport(body);
      if (res.status !== "success") throw new Error(res.message || "导入失败");
      toast.success(`导入完成：新增 ${res.summary?.created}，重复 ${res.summary?.duplicates}，错误 ${res.summary?.errors}`);
    } catch (e: any) {
      toast.error(e?.response?.data?.message || e.message || "导入失败");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-4">
      <ConfirmTable data={data} onChange={setData} />
      <div className="flex justify-end gap-2">
        <Button variant="outline" onClick={() => setData(initial)} disabled={loading}>重置</Button>
        <Button onClick={handleConfirm} disabled={loading}>{loading ? "导入中..." : "确认导入"}</Button>
      </div>
    </div>
  );
}


