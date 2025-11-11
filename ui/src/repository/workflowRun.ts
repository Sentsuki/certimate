import { type RecordSubscription } from "pocketbase";

import { type WorkflowRunModel } from "@/domain/workflowRun";

import { COLLECTION_NAME_WORKFLOW_RUN, getPocketBase } from "./_pocketbase";

const _commonFields = ["id", "status", "trigger", "startedAt", "endedAt", "error", "created", "updated", "deleted"];
const _expandFields = ["expand.workflowRef.id", "expand.workflowRef.name", "expand.workflowRef.description"];

export const list = async ({
  workflowId,
  page = 1,
  perPage = 10,
  expand = false,
}: {
  workflowId?: string;
  page?: number;
  perPage?: number;
  expand?: boolean;
}) => {
  const pb = getPocketBase();

  const filters: string[] = [];
  if (workflowId) {
    filters.push(pb.filter("workflowRef={:workflowId}", { workflowId: workflowId }));
  }

  return await pb.collection(COLLECTION_NAME_WORKFLOW_RUN).getList<WorkflowRunModel>(page, perPage, {
    expand: expand ? ["workflowRef"].join(",") : void 0,
    fields: [..._commonFields, ..._expandFields].join(","),
    filter: filters.join(" && "),
    sort: "-created",
    requestKey: null,
  });
};

export const get = async (id: string) => {
  return await getPocketBase()
    .collection(COLLECTION_NAME_WORKFLOW_RUN)
    .getOne<WorkflowRunModel>(id, {
      expand: ["workflowRef"].join(","),
      fields: ["*", ..._expandFields].join(","),
      requestKey: null,
    });
};

export const remove = async (record: MaybeModelRecordWithId<WorkflowRunModel> | MaybeModelRecordWithId<WorkflowRunModel>[]) => {
  const pb = getPocketBase();

  if (Array.isArray(record)) {
    const batch = pb.createBatch();
    for (const item of record) {
      batch.collection(COLLECTION_NAME_WORKFLOW_RUN).delete(item.id);
    }
    const res = await batch.send();
    return res.every((e) => e.status >= 200 && e.status < 400);
  } else {
    await pb.collection(COLLECTION_NAME_WORKFLOW_RUN).delete(record.id!);
    return true;
  }
};

export const subscribe = async (id: string, cb: (e: RecordSubscription<WorkflowRunModel>) => void) => {
  return getPocketBase().collection(COLLECTION_NAME_WORKFLOW_RUN).subscribe(id, cb);
};

export const unsubscribe = async (id: string) => {
  return getPocketBase().collection(COLLECTION_NAME_WORKFLOW_RUN).unsubscribe(id);
};
