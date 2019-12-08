<template>
	<div class="app-container">
		<el-breadcrumb separator-class="el-icon-arrow-right">
			<el-breadcrumb-item :to="{ path: '/welcome' }">首页</el-breadcrumb-item>
			<el-breadcrumb-item>系统管理</el-breadcrumb-item>
			<el-breadcrumb-item>菜单管理</el-breadcrumb-item>
		</el-breadcrumb>
		<div class="filter-container">
			<el-input v-model="listQuery.key" placeholder="请输入内容" clearable prefix-icon="el-icon-search" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter" @clear="handleFilter" />
			<el-select v-model="listQuery.type" placeholder="类型" clearable style="width: 90px" class="filter-item" @change="handleFilter">
				<el-option v-for="item in menuTypeOptions" :key="item.key" :label="item.display_name" :value="item.key" />
			</el-select>
			<SelectTree v-model="listQuery.parentId" class="filter-item" :props="propsSelectTree" :options="optionDataSelectTree" :value="valueIdSelectTree" :clearable="true" :accordion="true" @getValue="getSelectTreeValue($event, 1)" />
			<el-select v-model="listQuery.sort" style="width: 140px" class="filter-item" @change="handleFilter">
				<el-option v-for="item in sortOptions" :key="item.key" :label="item.label" :value="item.key" />
			</el-select>
			<el-button class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">
				{{ "搜索" }}
			</el-button>
			<el-button v-if="permissionList.add" class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-edit" @click="handleCreate">
				{{ "添加" }}
			</el-button>
			<el-button v-if="permissionList.del" class="filter-item" type="danger" icon="el-icon-delete" @click="handleBatchDel">
				{{ "删除" }}
			</el-button>
		</div>
		<!-- stripe属性可以创建带斑马纹的表格 -->
		<!-- border带边框 -->
		<!-- highlight-current-row属性即可实现单选 -->
		<!-- sort-change当表格的排序条件发生变化的时候会触发该事件 -->
		<!-- sort-change当表格的排序条件发生变化的时候会触发该事件 -->
		<!-- selection-change	当选择项发生变化时会触发该事件 -->
		<!-- v-loading来控制隐显 -->
		<el-table :key="tableKey" v-loading="listLoading" stripe :data="list" border fit highlight-current-row style="width: 100%;" @sort-change="sortChange" @selection-change="handleSelectionChange">
			<!-- 单选框 -->
			<el-table-column type="selection" width="55" />
			<el-table-column label="名称" align="center">
				<template slot-scope="scope">
					<span>{{ scope.row.name }}</span>
				</template>
			</el-table-column>
			<el-table-column label="代码" align="center">
				<template slot-scope="scope">
					<span>{{ scope.row.code }}</span>
				</template>
			</el-table-column>
			<el-table-column label="图标" align="center">
				<template slot-scope="scope">
					<span>{{ scope.row.icon }}</span>
				</template>
			</el-table-column>
			<el-table-column label="排序值" prop="sequence" sortable="custom" align="center">
				<template slot-scope="scope">
					<span>{{ scope.row.sequence }}</span>
				</template>
			</el-table-column>
			<el-table-column label="状态" prop="status" sortable="custom" align="center">
				<template slot-scope="scope">
					<span>{{ scope.row.status | statusFilter }}</span>
				</template>
			</el-table-column>
			<el-table-column label="操作" align="center" width="230" class-name="small-padding fixed-width">
				<template slot-scope="{ row }">
					<el-button v-if="permissionList.view" size="mini" type="success" @click="handleDetail(row.id)">
						{{ "查看" }}
					</el-button>
					<el-button v-if="permissionList.update" type="primary" size="mini" @click="handleUpdate(row.id)">
						{{ "编辑" }}
					</el-button>
					<el-button v-if="permissionList.del" size="mini" type="danger" @click="handleDelete(row)">
						{{ "删除" }}
					</el-button>
				</template>
			</el-table-column>
		</el-table>
		<!-- 生成底部的页码 -->
		<pagination v-show="total > 0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />
		<!-- visible.sync=dialogFormVisible,可视性，由点击按钮来切换 -->
		<el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
			<el-form ref="dataForm" v-loading="loading" element-loading-text="正在执行" element-loading-background="rgba(255,255,255,0.7)" :rules="rules" :model="temp" label-position="left" label-width="80px" style="width: 400px; margin-left:50px;">
				<el-form-item label="菜单类型" prop="menuType">
					<!-- v-model的值为当前被选中的el-option的 value 属性值 -->
					<el-select v-model.number="temp.menuType" type="number" placeholder="菜单类型">
						<el-option v-for="item in menuTypeOptions" :key="item.key" :label="item.display_name" :value="item.key" />
					</el-select>
				</el-form-item>
				<el-form-item label="操作类型" prop="operateType">
					<el-select v-model.number="temp.operateType" placeholder="操作类型">
						<el-option v-for="item in menuOperateTypeOptions" :key="item.key" :label="item.display_name" :value="item.key" />
					</el-select>
				</el-form-item>
				<el-form-item label="父级" prop="parentId">
					<SelectTree v-model.="temp.parentId" type="number" :props="propsSelectTree" :options="optionDataSelectTree" :value="valueIdSelectTree2" :clearable="true" :accordion="true" @getValue="getSelectTreeValue($event, 2)" />
				</el-form-item>
				<el-form-item label="名称" prop="name">
					<el-input v-model="temp.name" />
				</el-form-item>
				<el-form-item label="菜单url" prop="url">
					<el-input v-model="temp.url" />
				</el-form-item>
				<el-form-item label="代码" prop="code">
					<el-input v-model="temp.code" />
				</el-form-item>
				<el-form-item label="图标" prop="icon">
					<el-input v-model="temp.icon" />
				</el-form-item>
				<el-form-item label="排序值" prop="sequence">
					<el-input v-model.number="temp.sequence" type="number" />
				</el-form-item>
				<el-form-item label="状态" prop="status">
					<el-radio-group v-model.number="temp.status" type="number">
						<el-radio :label="1">启用</el-radio>
						<el-radio :label="2">未启用</el-radio>
					</el-radio-group>
				</el-form-item>
				<el-form-item label="备注" prop="memo">
					<el-input v-model="temp.memo" />
				</el-form-item>
			</el-form>
			<div v-if="
          dialogStatus !== 'detail' ? (loading === true ? false : true) : false
        " slot="footer" class="dialog-footer">
				<el-button @click="dialogFormVisible = false">
					{{ "取消" }}
				</el-button>
				<el-button type="primary" @click="dialogStatus === 'create' ? createData() : updateData()">
					{{ "确定" }}
				</el-button>
			</div>
		</el-dialog>
	</div>
</template>

<script>
import { requestList, requestDetail, requestUpdate, requestCreate, requestAll, requestDelete, requestMenuButton } from '@/api/sys/menu'

// 下面的分页组件
import Pagination from '@/components/Pagination' // Secondary package based on el-pagination
// 选择器组件，便于搜索查询
import SelectTree from '@/components/TreeSelect'
import { checkAuthAdd, checkAuthDel, checkAuthView, checkAuthUpdate } from '@/utils/permission'
// 类型 ，便于搜索
const menuTypeOptions = [
    { key: 1, display_name: '模块' },
    { key: 2, display_name: '菜单' },
    { key: 3, display_name: '操作' }
]
// 操作类型，用于编辑页面来声明这个api的作用
const menuOperateTypeOptions = [
    { key: 'none', display_name: '非操作类型' },
    { key: 'add', display_name: '新增' },
    { key: 'del', display_name: '删除' },
    { key: 'view', display_name: '查看' },
    { key: 'update', display_name: '编辑' },
    { key: 'list', display_name: '分页api' },
    { key: 'setrolemenu', display_name: '分配角色菜单权限' },
    { key: 'setadminrole', display_name: '分配管理员角色' }
]

export default {
    name: 'Menu',
    components: { Pagination, SelectTree },

    filters: {
        // 过滤器，1就显示的是启用，2就是不启用
        statusFilter(status) {
            const statusMap = {
                1: '启用',
                2: '不启用'
            }
            return statusMap[status]
        }
    },
    data() {
        return {
            valueIdSelectTree: 0,
            valueIdSelectTree2: 0,
            propsSelectTree: {
                value: 'id',
                label: 'name',
                children: 'children',
                placeholder: '父级'
            },
            propsSelectlist: [],
            propsSelectlist2: [
                { id: 0, parentId: -1, name: '顶级' }
            ],
            operationList: [],
            permissionList: {
                add: false,
                del: false,
                view: false,
                update: false
            },
            tableKey: 0,
            list: [],
            total: 0,
            listLoading: true,
            loading: true,
            listQuery: {
                page: 1,
                limit: 20,
                parentId: undefined,
                key: undefined,
                menuType: undefined,
                sort: '-id'
            },
            menuTypeOptions,
            menuOperateTypeOptions,
            sortOptions: [
                { label: 'ID Ascending', key: '+id' },
                { label: 'ID Descending', key: '-id' },
                { label: 'sequence Ascending', key: '+sequence' },
                { label: 'sequence Descending', key: '-sequence' },
                { label: 'status Ascending', key: '+status' },
                { label: 'status Descending', key: '-status' }
            ],
            temp: {
                id: 0,
                memo: '',
                name: '',
                url: '',
                code: '',
                icon: 'list',
                parentId: 0,
                menuType: 2,
                status: 1,
                sequence: 10
            },
            dialogFormVisible: false,
            dialogStatus: '',
            textMap: {
                update: '编辑',
                create: '添加',
                detail: '详情'
            },
            rules: {
                type: [{ required: true, message: '请选择类型', trigger: 'change' }],
                name: [{ required: true, message: '请输入菜单名称', trigger: 'blur' }],
                code: [{ required: true, message: '请输入菜单代码', trigger: 'blur' }],
                sequence: [{ required: true, message: '请输入排序值', trigger: 'blur' }]
            },
            multipleSelection: []
        }
    },
    computed: {
        optionDataSelectTree() {
            const cloneData = JSON.parse(JSON.stringify(this.propsSelectlist))
            return cloneData.filter(father => {
                const branchArr = cloneData.filter(child => father.id === child.parentId)
                branchArr.length > 0 ? father.children = branchArr : ''
                return father.parentId === this.propsSelectlist[0].parentId
            })
        }
    },
    created() {
        this.getMenuButton()
        this.getList()
        this.getAll()
    },
    methods: {
        checkPermission() {
            this.permissionList.add = checkAuthAdd(this.operationList)
            this.permissionList.del = checkAuthDel(this.operationList)
            this.permissionList.view = checkAuthView(this.operationList)
            this.permissionList.update = checkAuthUpdate(this.operationList)
        },
        getMenuButton() {
            requestMenuButton('Menu').then(response => {
                this.operationList = response.data
            }).then(() => {
                this.checkPermission()
            })
        },
        getList() {
            this.listLoading = true
            requestList(this.listQuery).then(response => {
                console.log(response.data.items)
                this.list = response.data.items
                this.total = response.data.total
                this.listLoading = false
            })
        },
        getAll() {
            requestAll().then(response => {
                if (response.data) {
                    this.propsSelectlist = response.data
                } else {
                    this.propsSelectlist = this.propsSelectlist2
                }
            })
        },
        handleFilter() {
            this.listQuery.parentId = this.valueIdSelectTree
            this.listQuery.page = 1
            this.getList()
        },
        sortChange(data) {
            console.log(this.list)
            const { prop, order } = data
            if (order === 'ascending') {
                this.listQuery.sort = '+' + prop
            } else if (order === 'descending') {
                this.listQuery.sort = '-' + prop
            } else {
                this.listQuery.sort = undefined
            }
            this.handleFilter()
        },
        resetTemp() {
            this.valueIdSelectTree2 = 0
            this.temp = {
                id: 0,
                memo: '',
                name: '',
                url: '',
                code: '',
                icon: 'list',
                operateType: 'none',
                parentId: 0,
                menuType: 2,
                status: 1,
                sequence: 10
            }
        },
        handleCreate() {
            this.resetTemp()
            this.dialogStatus = 'create'
            this.dialogFormVisible = true
            this.loading = false
            this.$nextTick(() => {
                this.$refs['dataForm'].clearValidate()
            })
        },
        createData() {
            this.$refs['dataForm'].validate((valid) => {
                if (valid) {
                    this.loading = true
                    this.temp.parentId = this.valueIdSelectTree2
                    requestCreate(this.temp).then(response => {
                        this.temp.id = response.data.id
                        this.list.unshift(this.temp)
                        this.dialogFormVisible = false
                        this.$notify({
                            title: '成功',
                            message: '创建成功',
                            type: 'success',
                            duration: 2000
                        })
                        this.total = this.total + 1
                        this.getAll()
                    }).catch(() => {
                        this.loading = false
                    })
                }
            })
        },
        handleDetail(id) {
            this.loading = true
            requestDetail(id).then(response => {
                this.loading = false
                this.temp = response.data
                this.valueIdSelectTree2 = this.temp.parentId
            })
            this.dialogStatus = 'detail'
            this.dialogFormVisible = true
            this.$nextTick(() => {
                this.$refs['dataForm'].clearValidate()
            })
        },
        handleUpdate(id) {
            this.loading = true
            requestDetail(id).then(response => {
                this.loading = false
                this.temp = response.data
                this.valueIdSelectTree2 = this.temp.parentId
            })
            this.dialogStatus = 'update'
            this.dialogFormVisible = true
            this.$nextTick(() => {
                this.$refs['dataForm'].clearValidate()
            })
        },
        updateData() {
            this.$refs['dataForm'].validate((valid) => {
                if (valid) {
                    this.loading = true
                    this.temp.parentId = this.valueIdSelectTree2
                    const tempData = Object.assign({}, this.temp)
                    requestUpdate(tempData).then(() => {
                        for (const v of this.list) {
                            if (v.id === this.temp.id) {
                                const index = this.list.indexOf(v)
                                this.list.splice(index, 1, this.temp)
                                break
                            }
                        }
                        this.dialogFormVisible = false
                        this.$notify({
                            title: '成功',
                            message: '更新成功',
                            type: 'success',
                            duration: 2000
                        })
                        this.getAll()
                    }).catch(() => {
                        this.loading = false
                    })
                }
            })
        },
        handleDelete(row) {
            var ids = []
            ids.push(row.id)
            this.$confirm('是否确定删除?', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                requestDelete(ids).then(() => {
                    this.$message({
                        message: '删除成功',
                        type: 'success'
                    })
                    this.total = this.total - 1
                    const index = this.list.indexOf(row)
                    this.list.splice(index, 1)
                    this.getAll()
                })
            }).catch(() => {
                this.$message({
                    type: 'info',
                    message: '已取消删除'
                })
            })
        },
        getSelectTreeValue(value, type) {
            if (type === 1) {
                this.valueIdSelectTree = value
                this.handleFilter()
            } else {
                this.valueIdSelectTree2 = value
            }
        },
        handleSelectionChange(val) {
            this.multipleSelection = val
        },
        handleBatchDel() {
            if (this.multipleSelection.length === 0) {
                this.$message({
                    message: '未选中任何行',
                    type: 'warning',
                    duration: 2000
                })
                return
            }
            var ids = []
            for (const v of this.multipleSelection) {
                ids.push(v.id)
            }
            this.$confirm('是否确定删除?', '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                requestDelete(ids).then(() => {
                    this.$message({
                        message: '删除成功',
                        type: 'success'
                    })
                    for (const row of this.multipleSelection) {
                        this.total = this.total - 1
                        const index = this.list.indexOf(row)
                        this.list.splice(index, 1)
                    }
                    this.getAll()
                })
            }).catch(() => {
                this.$message({
                    type: 'info',
                    message: '已取消删除'
                })
            })
        }
    }
}
</script>
