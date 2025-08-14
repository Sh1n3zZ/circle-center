import { useCallback, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { toast } from "sonner";
import { xmlioApi } from "@/api/manager/xmlio";
import type { IconImportComponent } from "@/api/manager/types";

export default function ProjectDetailXmlImportPanel({ onParsed }: { onParsed: (components: IconImportComponent[]) => void }) {
  const [appfilter, setAppfilter] = useState<File | null>(null);
  const [appmap, setAppmap] = useState<File | null>(null);
  const [theme, setTheme] = useState<File | null>(null);
  const [loading, setLoading] = useState(false);

  const handleDrop = useCallback((e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    const files = Array.from(e.dataTransfer.files || []);
    files.forEach((f) => {
      const name = f.name.toLowerCase();
      if (name.includes("appfilter")) setAppfilter(f);
      else if (name.includes("appmap")) setAppmap(f);
      else if (name.includes("theme")) setTheme(f);
      else if (name.endsWith(".xml")) {
        // Fallback: try to guess by size or first matched slot
        if (!appfilter) setAppfilter(f);
        else if (!appmap) setAppmap(f);
        else if (!theme) setTheme(f);
      }
    });
  }, [appfilter, appmap, theme]);

  const parseNow = async () => {
    if (!appfilter && !appmap && !theme) {
      toast.error("请先选择或拖拽至少一个 XML 文件");
      return;
    }
    try {
      setLoading(true);
      const res = await xmlioApi.parseForm(appfilter ?? undefined, appmap ?? undefined, theme ?? undefined);
      if (res.status !== "success" || !res.components) {
        throw new Error(res.message || "解析失败");
      }
      onParsed(res.components);
      toast.success(`解析成功，共 ${res.components.length} 条`);
    } catch (e: any) {
      toast.error(e?.response?.data?.message || e.message || "解析失败");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-4">
      <div
        className="border-dashed border rounded-md p-6 text-center text-sm text-muted-foreground"
        onDragOver={(e) => e.preventDefault()}
        onDrop={handleDrop}
      >
        将 appfilter.xml / appmap.xml / theme_resources.xml 拖拽到此处，或使用下方选择
      </div>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
        <div className="space-y-1">
          <div className="text-xs">appfilter.xml</div>
          <Input type="file" accept=".xml" onChange={(e) => setAppfilter(e.target.files?.[0] || null)} />
          {appfilter ? <div className="text-xs text-muted-foreground">{appfilter.name}</div> : null}
        </div>
        <div className="space-y-1">
          <div className="text-xs">appmap.xml</div>
          <Input type="file" accept=".xml" onChange={(e) => setAppmap(e.target.files?.[0] || null)} />
          {appmap ? <div className="text-xs text-muted-foreground">{appmap.name}</div> : null}
        </div>
        <div className="space-y-1">
          <div className="text-xs">theme_resources.xml</div>
          <Input type="file" accept=".xml" onChange={(e) => setTheme(e.target.files?.[0] || null)} />
          {theme ? <div className="text-xs text-muted-foreground">{theme.name}</div> : null}
        </div>
      </div>
      <Separator />
      <div className="flex justify-end">
        <Button onClick={parseNow} disabled={loading}>
          {loading ? "解析中..." : "解析预览"}
        </Button>
      </div>
    </div>
  );
}


