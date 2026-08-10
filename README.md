# ssmctl

AWS SSM 세션을 `.ssm.yml` 설정으로 간편하게 실행하는 Go CLI.

```bash
ssmctl dev-lg-jenkins
```

## 설치

```bash
go install github.com/heebin2/ssmctl/cmd@latest
```

## 설정

프로젝트 루트에 `.ssm.yml` 생성:

```yaml
global:
  user: ec2-user

instances:
  dev-lg-jenkins:
    target: i-010566b1d2073afd3
  prod-app:
    target: i-0abcde1234567890f
    user: ubuntu
```

## 사용

```bash
ssmctl list              # 인스턴스 목록 출력
ssmctl dev-lg-jenkins   # dev-lg-jenkins에 접속
```

옵션:
- `-config path`: 설정 파일 경로 (기본값: `.ssm.yml`)

## 구현

- `cmd/main.go`: 진입점
- `internal/ssm/`: 설정, 세션 로직
