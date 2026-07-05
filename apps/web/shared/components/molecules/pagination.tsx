"use client";

import { Button } from "@/shared/components/atoms/button";
import { SelectField } from "@/shared/components/atoms/select-field";

const defaultPageSizeOptions = [5, 10, 20];

type PaginationProps = {
  page: number;
  pageSize: number;
  totalPages: number;
  totalItems: number;
  firstRecord: number;
  lastRecord: number;
  pageSizeOptions?: number[];
  onPageChange: (page: number) => void;
  onPageSizeChange: (pageSize: number) => void;
};

export function Pagination({
  page,
  pageSize,
  totalPages,
  totalItems,
  firstRecord,
  lastRecord,
  pageSizeOptions = defaultPageSizeOptions,
  onPageChange,
  onPageSizeChange,
}: PaginationProps) {
  return (
    <div className="overflow-x-auto border-t border-slate-100 bg-white px-5 py-3 text-sm text-muted-foreground">
      <div className="flex min-w-max items-center justify-between gap-6">
        <div className="flex items-center gap-3">
          <p className="font-medium text-slate-600">
            Showing {firstRecord}-{lastRecord} of {totalItems}
          </p>
          <SelectField
            className="h-9 w-28 shrink-0"
            value={String(pageSize)}
            options={pageSizeOptions.map((size) => ({ value: String(size), label: `${size} / page` }))}
            onChange={(event) => onPageSizeChange(Number(event.target.value))}
          />
        </div>

        <div className="ml-auto flex items-center gap-2">
          <Button className="shrink-0" type="button" variant="outline" size="sm" disabled={page <= 1} onClick={() => onPageChange(Math.max(1, page - 1))}>
            Prev
          </Button>
          <span className="shrink-0 px-2 text-slate-600">
            Page {page} of {totalPages}
          </span>
          <Button className="shrink-0" type="button" variant="outline" size="sm" disabled={page >= totalPages} onClick={() => onPageChange(Math.min(totalPages, page + 1))}>
            Next
          </Button>
        </div>
      </div>
    </div>
  );
}
