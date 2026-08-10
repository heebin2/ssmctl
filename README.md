# ssmctl

AWS SSM 세션을 `~/.ssm.yml` 설정으로 간편하게 실행하는 Go CLI.

```bash
ssmctl my-instance
```

## 설치

```bash
go install github.com/heebin2/ssmctl/cmd@latest
```

## 설정

홈 디렉터리에 `~/.ssm.yml` 생성:

```yaml
global:
  user: ec2-user

instances:
  my-instance:
    target: i-1234567890abcdef0
  another-instance:
    target: i-0abcdef1234567890
    user: ubuntu
```

## 사용

```bash
ssmctl list           # 인스턴스 목록 출력
ssmctl my-instance    # my-instance에 접속
```

옵션:
- `-config path`: 설정 파일 경로 (기본값: `~/.ssm.yml`)

## 구현

- `cmd/main.go`: 진입점
- `internal/ssm/`: 설정, 세션 로직
