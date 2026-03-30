import mlflow

mlflow.set_tracking_uri("http://127.0.0.1:5001")

with mlflow.start_run():
    mlflow.log_param("model_type", "random_forest")
    mlflow.log_param("n_estimators", 100)
    mlflow.log_metric("accuracy", 0.95)
    mlflow.log_metric("loss", 0.12)
    print("Run geloggt ✅")