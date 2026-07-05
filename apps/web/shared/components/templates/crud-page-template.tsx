"use client";

import { PageShell } from "@/shared/components/organisms/page-shell";

type CrudPageTemplateProps = {
  title: string;
  subtitle?: string;
  toolbar?: React.ReactNode;
  table: React.ReactNode;
  pagination?: React.ReactNode;
  modal?: React.ReactNode;
  dialogs?: React.ReactNode;
};

export function CrudPageTemplate({ title, subtitle, toolbar, table, pagination, modal, dialogs }: CrudPageTemplateProps) {
  return (
    <>
      <PageShell title={title} subtitle={subtitle} toolbar={toolbar}>
        {table}
        {pagination}
      </PageShell>
      {modal}
      {dialogs}
    </>
  );
}
