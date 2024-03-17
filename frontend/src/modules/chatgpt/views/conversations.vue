<template>
	<cl-crud ref="Crud">
		<cl-row>
			<!-- 刷新按钮 -->
			<cl-refresh-btn />
			<!-- 新增按钮 -->
			<cl-add-btn />
			<!-- 删除按钮 -->
			<cl-multi-delete-btn />
			<cl-flex1 />
			<!-- 关键字搜索 -->
			<cl-search-key />
		</cl-row>

		<cl-row>
			<!-- 数据表格 -->
			<cl-table ref="Table" />
		</cl-row>

		<cl-row>
			<cl-flex1 />
			<!-- 分页控件 -->
			<cl-pagination />
		</cl-row>

		<!-- 新增、编辑 -->
		<cl-upsert ref="Upsert" />
	</cl-crud>
</template>

<script lang="ts" name="chatgpt-conversations" setup>
import { useCrud, useTable, useUpsert } from "@cool-vue/crud";
import { useCool } from "/@/cool";

const { service } = useCool();

// cl-upsert 配置
const Upsert = useUpsert({
	items: [
		{ label: "用户token", prop: "usertoken", required: true, component: { name: "el-input" } },
		{ label: "会话id", prop: "convid", required: true, component: { name: "el-input" } },
		{ label: "会话标题", prop: "title", required: true, component: { name: "el-input" } },
		{ label: "官网账号邮箱", prop: "email", required: true, component: { name: "el-input" } }
	]
});

// cl-table 配置
const Table = useTable({
	columns: [
		{ type: "selection" },
		{ label: "id", prop: "id" },
		{ label: "创建时间", prop: "createTime" },
		{ label: "更新时间", prop: "updateTime" },
		{ label: "用户token", prop: "usertoken" },
		{ label: "会话id", prop: "convid" },
		{ label: "会话标题", prop: "title" },
		{ label: "官网账号邮箱", prop: "email" },
		{ label: "组织ID", prop: "chatgptaccountid" },
		{ type: "op", buttons: ["edit", "delete"] }
	]
});

// cl-crud 配置
const Crud = useCrud(
	{
		service: service.chatgpt.conversations
	},
	(app) => {
		app.refresh();
	}
);
</script>
