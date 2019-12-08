<!--  -->
<template>
	<div class="login_container">
		<div class="login_box">

			<el-form ref="loginFormRef" label-width="0px" class="login_form" :model="loginForm" :rules="loginFormRules">
				<el-form-item prop="username">
					<el-input v-model="loginForm.username" prefix-icon="iconfont icon-user"> </el-input>
				</el-form-item>
				<el-form-item prop="password">
					<el-input v-model="loginForm.password" prefix-icon="iconfont icon-3702mima" type="password"></el-input>
				</el-form-item>
				<el-form-item class="btns">
					<el-button type="primary" @click="login">登录</el-button>

					<el-button type="info" @click="resetLoginForm">重置</el-button>
				</el-form-item>
			</el-form>
		</div>

	</div>
</template>
<script>
export default {
	components: {},
	methods: {
		resetLoginForm() {
			console.log('11')
			this.$refs.loginFormRef.resetFields()
		},
		login() {
			this.$refs.loginFormRef.validate(valid => {
				if (valid) {
					this.$store
						.dispatch('user/login', this.loginForm)
						.then(() => {
							// this.$message.success('登录成功!')
							this.$router.push('/home')
						})
						.catch(() => {})
				}
			})
		}
	},
	data() {
		return {
			loginForm: {
				username: 'admin',
				password: '123456'
			},
			loginFormRules: {
				username: [
					{
						required: true,
						message: '请输入登录名称',
						trigger: 'blur'
					}
				],
				password: [
					{
						required: true,
						message: '请输入登录密码',
						trigger: 'blur'
					},
					{
						min: 3,
						max: 15,
						message: '长度在 3 到 15 个字符',
						trigger: 'blur'
					}
				]
			}
		}
	}
}
</script>
 
<style scoped lang = "less">
.login_container {
	background-color: #b5ade4;
	height: 100%;
}
.login_box {
	width: 450px;
	height: 300px;
	background-color: #fff;
	border-radius: 3px;
	position: absolute;
	left: 50%;
	top: 50%;
	transform: translate(-50%, -50%);
}

.btns {
	display: flex;
	justify-content: flex-end;
}
.login_form {
	position: absolute;
	width: 100%;
	bottom: 0;
	padding: 0 20px;
	box-sizing: border-box;
}
</style>