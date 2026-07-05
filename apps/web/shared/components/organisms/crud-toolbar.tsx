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
    <div className="flex flex-col gap-3 lg:flex-row lg:items-center">
      <SearchBox className="lg:w-72" value={search} placeholder={searchPlaceholder} onChange={onSearch} />

      <div className="flex flex-wrap gap-2">
        {filters}
        <Button type="button" variant="outline" onClick={onExportExcel}>
          <Download className="h-4 w-4" />
          Excel
        </Button>
        <Button type="button" variant="outline" onClick={onExportPdf}>
          <FileText className="h-4 w-4" />
          PDF
        </Button>
        <Button type="button" onClick={onAdd}>
          <Plus className="h-4 w-4" />
          {addLabel}
        </Button>
      </div>
    </div>
  );
}
