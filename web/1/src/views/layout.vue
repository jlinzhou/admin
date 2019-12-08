<!--  -->
<template>
	<!-- <div> -->

	<el-container class="home-container">
		<el-header>
			<div>
				<!-- <img src="@/assets/lutou.png" alt=""> -->
				<span>后台管理系统</span>
				<router-link to="/welcome">
					<div class="backHome">
						<i class="el-icon-s-home"></i>
					</div>
				</router-link>
			</div>
			<div>
				<el-dropdown>
					<span class="el-dropdown-link">
						操作<i class="el-icon-arrow-down el-icon--right"></i>
					</span>
					<el-dropdown-menu slot="dropdown">
						<el-dropdown-item>
							<router-link to="/welcome" style="text-decoration:none">
								<span>首页</span>
							</router-link>
						</el-dropdown-item>
						<el-dropdown-item>
							<span @click="handleEditPwd">修改密码</span>
						</el-dropdown-item>
						<el-dropdown-item>
							<span @click="logout">退出登录</span>
						</el-dropdown-item>
					</el-dropdown-menu>
				</el-dropdown>
			</div>
			<!-- <el-button type="info" @click="logout">退出</el-button> -->
		</el-header>
		<el-container>
			<!-- 第一个路由标题在meta里面 -->
			<!-- 子路由标题在儿子meta下面 -->
			<!-- 第一层判断 -->
			<el-aside :width="isCollapse ? '64px':'200px'">

				<div class="toggle-button" @click="toggleCollapse">|||</div>

				<el-menu background-color="#373d41" text-color="#fff" active-text-color="#409eff" :unique-opened="true" :default-active="activePath" :collapse="isCollapse" :collapse-transition="false" router>
					<el-submenu :index="item.name" v-for="item in menulist" :key="item.name">
						<template slot="title">
							<!-- <i :class="iconsObj[item.id]"></i> -->

							<span>{{item.meta.title}}</span>
						</template>
						<template v-for="subItem in item.children">
							<template v-for="subItemChi in subItem.children">
								<template v-if="!subItemChi.hidden">
									<el-menu-item :index="subItem.path" :key="subItemChi.name" @click="saveNavState('/'+subItem.path)">
										<i class="el-icon-menu"></i>
										<span> {{subItemChi.meta.title}}</span>
									</el-menu-item>
								</template>
							</template>
						</template>
					</el-submenu>

				</el-menu>
			</el-aside>
			<el-main>
				<router-view :key="$route.fullPath"> </router-view>
			</el-main>
		</el-container>
		<el-dialog title="修改密码" :visible.sync="dialogFormVisible">
			<el-form ref="dataForm" v-loading="loading" element-loading-text="正在执行" element-loading-background="rgba(255,255,255,0.7)" :rules="rules" :model="temp" label-position="left" label-width="120px" style="width: 400px; margin-left:50px;">
				<el-form-item label="原密码" prop="old_password">
					<el-input v-model="temp.old_password" show-password minlength="6" maxlength="20" />
				</el-form-item>
				<el-form-item label="新密码" prop="new_password">
					<el-input v-model="temp.new_password" placeholder="6-20位" show-password minlength="6" maxlength="20" />
				</el-form-item>
				<el-form-item label="再次输入新密码" prop="new_password_again">
					<el-input v-model="temp.new_password_again" placeholder="6-20位" show-password minlength="6" maxlength="20" />
				</el-form-item>
			</el-form>
			<div slot="footer" class="dialog-footer">
				<el-button @click="dialogFormVisible = false">
					{{ "取消" }}
				</el-button>
				<el-button type="primary" @click="editPwd()">
					{{ "确定" }}
				</el-button>
			</div>
		</el-dialog>
	</el-container>
	<!-- </div> -->

</template>

<script>
import { requestEditPwd } from '@/api/sys/user'
export default {
	name: 'layout',
	components: {},
	data() {
		return {
			menulist: [],
			iconsObj: {
				'125': 'iconfont icon-user',
				'103': 'iconfont icon-tijikongjian',
				'101': 'iconfont icon-shangpin',
				'102': 'iconfont icon-danju',
				'145': 'iconfont icon-baobiao'
			},
			isCollapse: false,
			activePath: '',
			dialogFormVisible: false,
			loading: true,
			temp: {
				old_password: '',
				new_password: '',
				new_password_again: ''
			},
			rules: {
				old_password: [
					{ required: true, message: '请输入旧密码', trigger: 'blur' }
				],
				new_password: [
					{
						min: 6,
						max: 20,
						required: true,
						message: '长度在 6 到 20 个字符',
						trigger: 'blur'
					}
				],
				new_password_again: [
					{
						min: 6,
						max: 20,
						required: true,
						message: '长度在 6 到 20 个字符',
						trigger: 'blur'
					}
				]
			}
		}
	},
	created() {
		this.getMenuList()
		this.activePath = window.sessionStorage.getItem('activePath')
	},
	methods: {
		// logout() {

		// 	removeToken()
		// 	this.$router.push('/login')
		// },
		async logout() {
			await this.$store.dispatch('user/logout')
			this.$router.push('/login')
		},
		//this.$store.getters.permission_routes
		getMenuList() {
			//this.menulist = this.$store.getters.permission_routes
			var allRoutes = this.$store.getters.permission_routes
			for (var i = 0; i < allRoutes.length; i++) {
				if (allRoutes[i].path !== '/home') {
					if (allRoutes[i].children) {
						this.menulist.push(allRoutes[i])
					}
				}
			}
		},

		toggleCollapse() {
			this.isCollapse = !this.isCollapse
		},
		saveNavState(activePath) {
			window.sessionStorage.setItem('activePath', activePath)
			this.activePath = activePath
		},
		resetTemp() {
			this.temp = {
				old_password: '',
				new_password: '',
				new_password_again: ''
			}
		},
		handleEditPwd() {
			this.resetTemp()
			this.dialogFormVisible = true
			this.loading = false
			this.$nextTick(() => {
				this.$refs['dataForm'].clearValidate()
			})
		},
		editPwd() {
			this.$refs['dataForm'].validate(valid => {
				if (valid) {
					if (
						this.temp.new_password !== this.temp.new_password_again
					) {
						this.$message.error('两次输入的密码不一致')
						return
					}
					this.loading = true
					requestEditPwd(this.temp)
						.then(response => {
							this.dialogFormVisible = false
							this.$notify({
								title: '成功',
								message: '修改成功',
								type: 'success',
								duration: 2000
							})
						})
						.catch(() => {
							this.loading = false
						})
				}
			})
		}
	}
}
</script>
 
<style scoped lang = "less">
.home-container {
	height: 100%;
}
.el-header {
	background-color: #373d41;
	display: flex;
	justify-content: space-between;
	padding-left: 0px;
	align-items: center;
	color: #fff;
	font-size: 20px;
	> div {
		display: flex;
		align-items: center;
		span {
			margin-left: 15px;
		}
	}
}
.el-aside {
	background-color: #333744;
	.el-menu {
		border-right: none;
	}
}
.el-main {
	background-color: #fcfcfc;
}
.iconfont {
	margin-right: 10px;
}
.toggle-button {
	background-color: #4a5064;
	font-size: 10px;
	line-height: 24px;
	color: #fff;
	text-align: center;
	letter-spacing: 0.2em;
	cursor: pointer;
}
.el-dropdown-link {
	cursor: pointer;
	color: #409eff;
}
.el-icon-arrow-down {
	font-size: 12px;
}
.backHome {
	text-align: center;
	font-size: 18px;
	color: rgb(255, 255, 255);
}
</style>