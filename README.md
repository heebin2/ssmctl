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

### 자동 초기화 (AWS EC2 인스턴스 자동 로드)

AWS 계정에 설정된 EC2 인스턴스를 자동으로 로드:

```bash
ssmctl init
```

`~/.ssm.yml`이 자동으로 생성되며, 실행 중인 EC2 인스턴스가 `Name` 태그로 등록됩니다.

### 수동 설정

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
ssmctl init           # 설정 초기화
ssmctl list           # 인스턴스 목록 출력
ssmctl my-instance    # my-instance에 접속
```

옵션:
- `-config path`: 설정 파일 경로 (기본값: `~/.ssm.yml`)

### zsh 자동완성

다음 명령을 한 번 실행해 자동완성을 활성화합니다.

```zsh
echo 'eval "$(ssmctl completion zsh)"' >> ~/.zshrc
source ~/.zshrc
```

이후 `ssmctl <TAB>`을 누르면 명령과 `~/.ssm.yml`의 인스턴스 이름이 표시됩니다.
`ssmctl -config /path/to/config.yml <TAB>`처럼 별도 설정 파일도 사용할 수 있습니다.

## 구현

- `cmd/main.go`: 진입점
- `internal/ssm/`: 설정, 세션, 초기화 로직
