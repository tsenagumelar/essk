"use client";

import { Download, FileText, Plus } from "lucide-react";
import { Button } from "@/shared/components/atoms/button";
import { SearchBox } from "@/shared/components/molecules/search-box";

type CrudToolbarProps = {
  search: string;
  searchPlaceholder?: string;
  filters?: React.ReactNode;
  addLabel?: string;
  onSearch: (value: string) => void;
  onAdd: () => void;
  onExportExcel: () => void;
  onExportPdf: () => void;
};

export function CrudToolbar({
  search,
  searchPlaceholder = "Search records",
  filters,
  addLabel = "Add",
  onSearch,
  onAdd,
  onExportExcel,
  onExportPdf,
}: CrudToolbarProps) {
  return (
    <div className="overflow-x-auto rounded-xl border border-slate-200 bg-slate-50/80 p-3">
      <div className="flex min-w-max items-center justify-between gap-6">
        <SearchBox
          className="w-[520px] shrink-0"
          value={search}
          placeholder={searchPlaceholder}
          onChange={onSearch}
        />

        <div className="ml-auto flex shrink-0 items-center gap-2">
          {filters}
          <Button
            className="shrink-0"
            type="button"
            variant="outline"
            onClick={onExportExcel}
          >
            <Download className="h-4 w-4" />
            Excel
          </Button>
          <Button
            className="shrink-0"
            type="button"
            variant="outline"
            onClick={onExportPdf}
          >
            <FileText className="h-4 w-4" />
            PDF
          </Button>
          <Button className="shrink-0" type="button" onClick={onAdd}>
            <Plus className="h-4 w-4" />
            {addLabel}
          </Button>
        </div>
      </div>
    </div>
  );
}
