
<template>
    <div>
        <el-breadcrumb separator-class="el-icon-arrow-right">
            <el-breadcrumb-item :to="{ path: '/home' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item>用户管理</el-breadcrumb-item>
            <el-breadcrumb-item>用户列表</el-breadcrumb-item>
        </el-breadcrumb>
        <el-card class="box-card">
            <el-row :gutter="20">
                <el-col :span="7">
                    <el-input placeholder="请输入内容" clearable v-model="listQuery.query" @clear="getList">
                        <el-button slot="append" icon="el-icon-search" @click="getList"></el-button>
                    </el-input>
                </el-col>
                <el-col :span="4">
                    <el-button type="primary" @click="handleCreate "> {{ "添加" }}</el-button>
                </el-col>
                <el-col :span="4">
                    <el-button type="danger">{{ "删除" }}</el-button>
                </el-col>
            </el-row>
            <el-table :data="list" border stripe>
                <el-table-column type="index">
                </el-table-column>
                <el-table-column prop="username" label="姓名">
                </el-table-column>
                <el-table-column prop="email" label="邮箱">
                </el-table-column>
                <el-table-column prop="mobile" label="电话">
                </el-table-column>
                <el-table-column prop="role_name" label="角色">
                </el-table-column>
                <el-table-column prop="mg_state" label="状态">
                    <template slot-scope="scope">
                        <el-switch v-model="scope.row.mg_state" @change="userStateChanged(scope.row)">
                        </el-switch>
                    </template>
                </el-table-column>
                <el-table-column label="操作" width="180px">
                    <template slot-scope="scope">

                        <el-button type="primary" icon="el-icon-edit" size="mini" @click="showEditDialog(scope.row.id)"></el-button>
                        <el-button type="danger" icon="el-icon-delete" size="mini"></el-button>
                        <el-tooltip class="dark" effect="dark" content="分配角色" :enterable="false" placement="top">
                            <el-button type="warning" icon="el-icon-setting" size="mini"></el-button>
                        </el-tooltip>
                    </template>
                </el-table-column>
            </el-table>

        </el-card>
        <pagination v-show="total > 0" :total="total" :page.sync="listQuery.pagenum" :limit.sync="listQuery.pagesize" @pagination="getList" />

        <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
            <!-- :rules="rules"  -->
            <el-form :model="temp" :rules="rules" ref="dataForm" label-width="70px" class="demo-ruleForm">
                <el-form-item label="用户名" prop="username">
                    <el-input v-model="temp.username"></el-input>
                </el-form-item>
                <el-form-item label="密码" prop="password">
                    <el-input v-model="temp.password"></el-input>
                </el-form-item>
                <el-form-item label="邮箱" prop="email">
                    <el-input v-model="temp.email"></el-input>
                </el-form-item>
                <el-form-item label="手机" prop="mobile">
                    <el-input v-model="temp.mobile"></el-input>
                </el-form-item>
            </el-form>
            <div v-if="
          dialogStatus !== 'detail' 
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

import Pagination from '@/components/Pagination'
import { requestList, requestDetail, requestCreate } from '@/api/user'
export default {
    components: { Pagination },

    data() {
        var checkEmail = (rule, value, cb) => {
            const regEmail = /^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$/
            if (regEmail.test(value)) {
                return cb()
            }
            cb(new Error('请输入合法的邮箱'))
        }
        var checkMobile = (rule, value, cb) => {
            const regMobile = /^1([38][0-9]|4[579]|5[0-3,5-9]|6[6]|7[0135678]|9[89])\d{8}$/
            if (regMobile.test(value)) {
                return cb()
            }
            cb(new Error('请输入合法的手机号'))
        }
        return {
            // 获取操作权限列表
            operationList: [],
            // 控制按钮的显示
            permissionList: {
                add: false,
                del: false,
                view: false,
                update: false
            },

            tableKey: 0,
            // 所有的数据
            list: [],
            // total为数据的总数
            total: 0,
            // 分页的参数
            listQuery: {
                query: '',
                pagenum: 1,
                pagesize: 2
                // page: 1,
                // limit: 20,
                // key: undefined,
                // status: undefined,
                // sort: '-exchangeId'
            },

            temp: {
                username: '',
                password: '',
                email: '',
                mobile: ''
            },
            // 对话框是否显示
            dialogFormVisible: false,
            // 是哪种对话框
            dialogStatus: '',
            // 对话框的标题
            textMap: {
                update: '编辑',
                create: '添加',
                detail: '详情'
            },
            rules: {
                username: [
                    {
                        required: true,
                        message: '请输入用户名',
                        trigger: 'blur'
                    },
                    {
                        min: 3,
                        max: 10,
                        message: '用户名的长度在3~10个之间',
                        trigger: 'blur'
                    }
                ],
                password: [
                    {
                        required: true,
                        message: '请输入密码',
                        trigger: 'blur'
                    },
                    {
                        min: 6,
                        max: 15,
                        message: '密码的长度在3~10个之间',
                        trigger: 'blur'
                    }
                ],
                email: [
                    {
                        required: true,
                        message: '请输入邮箱',
                        trigger: 'blur'
                    },
                    {
                        validator: checkEmail,
                        trigger: 'blur'
                    }
                ],
                mobile: [
                    {
                        required: true,
                        message: '请输入手机号',
                        trigger: 'blur'
                    },
                    {
                        validator: checkMobile,
                        trigger: 'blur'
                    }
                ]
            },

            multipleSelection: []

        }
    },
    // 响应后执行下面的函数
    created() {
        this.getList()
    },
    methods: {
        // 获取分页数据
        getList() {
            requestList(this.listQuery).then(response => {
                var res = response.data
                if (res.meta.status !== 200) {
                    return this.$message.error('获取用户列表失败!')
                }
                this.list = res.data.users
                this.total = res.data.total
                console.log(res)
            })
        },

        resetTemp() {
            this.temp = {
                username: '',
                password: '',
                email: '',
                mobile: ''
            }
        },
        // 点击创建后触发
        handleCreate() {
            this.resetTemp()
            this.dialogStatus = 'create'
            this.dialogFormVisible = true
            // $nextTick将回调延迟到下次 DOM 更新循环之后执行。
            this.$nextTick(() => {
                // clearValidate移除表单的校验结果
                this.$refs['dataForm'].clearValidate()
            })
        },
        createData() {
            // validate对整个表单进行校验的方法
            this.$refs['dataForm'].validate((valid) => {
                if (valid) {
                    this.loading = true
                    requestCreate(this.temp).then(response => {
                        this.temp.symbolProduct = response.data.symbolProduct
                        // unshift() 方法可向数组的开头添加一个或更多元素
                        this.list.unshift(this.temp)
                        this.dialogFormVisible = false
                        // 貌似是vue自带的提示功能
                        this.$notify({
                            title: '成功',
                            message: '创建成功',
                            type: 'success',
                            duration: 2000
                        })
                        this.total = this.total + 1
                    }).catch(() => {
                        this.loading = false
                    })
                }
            })
        },

        // 这里传入主键
        handleDetail(id) {
            this.loading = true
            requestDetail(id).then(response => {
                this.temp = response.data
            })
            this.dialogStatus = 'detail'
            this.dialogFormVisible = true
            this.$nextTick(() => {
                this.$refs['dataForm'].clearValidate()
            })
        },
        // 这里传入主键
        handleUpdate(id) {
            requestDetail(id).then(response => {
                this.temp = response.data
            })
            this.dialogStatus = 'update'
            this.dialogFormVisible = true
            this.$nextTick(() => {
                this.$refs['dataForm'].clearValidate()
            })
        },
        updateData() {
            // this.$refs['dataForm'].validate((valid) => {
            //     if (valid) {
            //         this.loading = true
            //         const tempData = Object.assign({}, this.temp)
            //         requestUpdate(tempData).then(() => {
            //             for (const v of this.list) {
            //                 if (v.exchangeId == this.temp.exchangeId && v.symbolProduct == this.temp.symbolProduct) {
            //                     const index = this.list.indexOf(v)
            //                     this.list.splice(index, 1, this.temp)
            //                     break
            //                 }
            //             }
            //             this.dialogFormVisible = false
            //             this.$notify({
            //                 title: '成功',
            //                 message: '更新成功',
            //                 type: 'success',
            //                 duration: 2000
            //             })
            //         }).catch(() => {
            //             this.loading = false
            //         })
            //     }
            // })
        },
        handleDelete(row) {
            //     var ids = []
            //     var data = { exchangeId: row.exchangeId, symbolProduct: row.symbolProduct }
            //     this.$confirm('是否确定删除?', '提示', {
            //         confirmButtonText: '确定',
            //         cancelButtonText: '取消',
            //         type: 'warning'
            //     }).then(() => {
            //         requestDelete(JSON.stringify({ ids: ids })).then(() => {
            //             this.$message({
            //                 message: '删除成功',
            //                 type: 'success'
            //             })
            //             this.total = this.total - 1
            //             const index = this.list.indexOf(row)
            //             this.list.splice(index, 1)
            //         })
            //     }).catch(() => {
            //         this.$message({
            //             type: 'info',
            //             message: '已取消删除'
            //         })
            //     })
        }
    }
}
</script>
