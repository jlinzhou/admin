<!--  -->
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
                    <el-input placeholder="请输入内容" clearable v-model="queryInfo.query" @clear="getUserList">
                        <el-button slot="append" icon="el-icon-search" @click="getUserList"></el-button>
                    </el-input>
                </el-col>
                <el-col :span="4">
                    <el-button type="primary" @click="addDialogVisible = true"> 添加用户</el-button>
                </el-col>
            </el-row>

            <el-table :data="userlist" border stripe>
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
            <el-pagination @size-change="handleSizeChange" @current-change="handleCurrentChange" :current-page="queryInfo.pagenum" :page-sizes="[1, 2, 5, 10]" :page-size="queryInfo.pagesize" layout="total, sizes, prev, pager, next, jumper" :total="total">
            </el-pagination>
        </el-card>
        <el-dialog title="添加用户" :visible.sync="addDialogVisible" width="50%" @close="addDialogClosed">
            <el-form :model="addForm" :rules="addFormRules" ref="addFormRef" label-width="70px" class="demo-ruleForm">
                <el-form-item label="用户名" prop="username">
                    <el-input v-model="addForm.username"></el-input>
                </el-form-item>
                <el-form-item label="密码" prop="password">
                    <el-input v-model="addForm.password"></el-input>
                </el-form-item>
                <el-form-item label="邮箱" prop="email">
                    <el-input v-model="addForm.email"></el-input>
                </el-form-item>
                <el-form-item label="手机" prop="mobile">
                    <el-input v-model="addForm.mobile"></el-input>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="addDialogVisible = false">取 消</el-button>
                <el-button type="primary" @click="addUser">确 定</el-button>
            </span>
        </el-dialog>
        <el-dialog title="修改用户" :visible.sync="editDialogVisible" width="50%">
            <span>这是一段信息</span>
            <span slot="footer" class="dialog-footer">
                <el-button @click="editDialogVisible = false">取 消</el-button>
                <el-button type="primary" @click="editDialogVisible = false">确 定</el-button>
            </span>
        </el-dialog>
    </div>
</template>

<script>
import { requestList, requestDetail, requestCreate } from '@/api/user'
export default {
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
            queryInfo: {
                query: '',
                pagenum: 1,
                pagesize: 2
            },
            userlist: [],
            total: 0,
            //控制
            addDialogVisible: false,
            editDialogVisible: false,
            addForm: {
                username: '',
                password: '',
                email: '',
                mobile: ''
            },
            editForm: {},
            addFormRules: {
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
            }
        }
    },
    components: {},
    created() {
        // this.getUserList2()
        this.getUserList()
    },
    methods: {
        // async getUserList() {
        // 	const { data: res } = await this.$http.get('users', {
        // 		params: this.queryInfo
        // 	})
        // 	if (res.meta.status !== 200) {
        // 		return this.$message.error('获取用户列表失败!')
        // 	}
        // 	this.userlist = res.data.users
        // 	this.total = res.data.total
        // 	console.log(res)
        // },
        getUserList() {
            requestList(this.queryInfo).then(response => {
                var res = response.data
                if (res.meta.status !== 200) {
                    return this.$message.error('获取用户列表失败!')
                }
                this.userlist = res.data.users
                this.total = res.data.total
                console.log(res)
            })
        },
        handleSizeChange(newSize) {
            this.queryInfo.pagesize = newSize
            this.getUserList()
        },
        handleCurrentChange(newPage) {
            this.queryInfo.pagenum = newPage
            this.getUserList()
        },
        async userStateChanged(userinfo) {
            const { data: res } = await this.$http.put(
                `users/${userinfo.id}/state/${userinfo.mg_state}`
            )
            if (res.meta.status !== 200) {
                userinfo.mg_state = !userinfo.mg_state
                return this.$message.error('更新用户失败!')
            }
            this.$message.success('更新用户状态成功!')
        },
        addDialogClosed() {
            this.$refs.addFormRef.resetFields()
        },
        addUser() {
            requestCreate(this.addForm).then(response => {
                if (res.meta.status !== 200) {
                    this.$message.error('添加用户失败!')
                }
                this.$message.success('添加用户成功!')
                this.addDialogVisible = false
                this.getUserList()
            })
            // this.$refs.addFormRef.validate(async valid => {
            //     if (!valid) return
            //     const { data: res } = await this.$http.post(
            //         'users',
            //         this.addForm
            //     )
            //     if (res.meta.status !== 200) {
            //         this.$message.error('添加用户失败!')
            //     }
            //     this.$message.success('添加用户成功!')
            //     this.addDialogVisible = false
            //     this.getUserList()
            // })
        },
        showEditDialog(id) {
            requestDetail(id).then(response => {
                var res = response.data
                if (res.meta.status !== 200) {
                    return this.$message.error('获取用户列表失败!')
                }
                this.editForm = res.data
                this.editDialogVisible = true
            })
        }
        // async showEditDialog(id) {
        //     console.log(id)
        //     const { data: res } = await this.$http.get('users/' + id)
        //     console.log(res)
        //     if (res.meta.status !== 200) {
        //         return this.$message.error('查询用户信息失败!')
        //     }
        //     this.editForm = res.data
        //     this.editDialogVisible = true
        // }
    }
}
</script>
 
<style scoped lang = "less">
</style>