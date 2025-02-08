import React, { useState } from "react";
import { Form, Input, Button, Typography, Card, message } from "antd";
import axios from "axios";
import { useNavigate } from "react-router-dom";

const { Title } = Typography;

const AuthPage = ({ type, onAuthSuccess }) => {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const isLogin = type === "login";

  const onFinish = async (values) => {
    setLoading(true);
    try {
      const url = `http://localhost:5000/${isLogin ? "login" : "signup"}`;
      await axios.post(url, values, {
        withCredentials: true
      });
      
      onAuthSuccess();
      navigate("/dashboard");
      message.success(isLogin ? "Вход выполнен!" : "Регистрация успешна!");
    } catch (error) {
      message.error(
        error.response?.data?.message || 
        (isLogin ? "Ошибка входа!" : "Ошибка регистрации!")
      );
    }
    setLoading(false);
  };

  return (
    <div style={{ 
      maxWidth: 400, 
      margin: "50px auto",
      padding: 20
    }}>
      <Card>
        <Title level={2} style={{ textAlign: "center" }}>
          {isLogin ? "Вход" : "Регистрация"}
        </Title>
        
        <Form
          name="auth-form"
          initialValues={{ remember: true }}
          onFinish={onFinish}
          layout="vertical"
        >
          <Form.Item
            label="Имя пользователя"
            name="username"
            rules={[
              { required: true, message: "Введите имя пользователя!" }
            ]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            label="Пароль"
            name="password"
            rules={[{ required: true, message: "Введите пароль!" }]}
          >
            <Input.Password />
          </Form.Item>

          <Form.Item>
            <Button 
              type="primary" 
              htmlType="submit" 
              loading={loading}
              block
            >
              {isLogin ? "Войти" : "Зарегистрироваться"}
            </Button>
          </Form.Item>

          <div style={{ textAlign: "center" }}>
            {isLogin ? (
              <span>
                Нет аккаунта?{" "}
                <Button type="link" onClick={() => navigate("/signup")}>
                  Зарегистрироваться
                </Button>
              </span>
            ) : (
              <span>
                Уже есть аккаунт?{" "}
                <Button type="link" onClick={() => navigate("/login")}>
                  Войти
                </Button>
              </span>
            )}
          </div>
        </Form>
      </Card>
    </div>
  );
};

export default AuthPage;