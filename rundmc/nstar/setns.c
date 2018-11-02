#include <stdio.h>
#include <sys/param.h>
#include <linux/sched.h>
#include <linux/fcntl.h>

/* nothing seems to define this... */
int setns(int fd, int nstype);

int nsenter(int pid) {
  pid_t tpid;
  char mntnspath[PATH_MAX];
  char usrnspath[PATH_MAX];
  int mntnsfd;
  int usrnsfd;

  tpid = (pid_t) pid;

  if(snprintf(mntnspath, sizeof(mntnspath), "/proc/%u/ns/mnt", tpid) == -1) {
    perror("snprintf ns mnt path");
    return 1;
  }

  mntnsfd = open(mntnspath, O_RDONLY);
  if(mntnsfd == -1) {
    perror("open mnt namespace");
    return 1;
  }

  if(snprintf(usrnspath, sizeof(usrnspath), "/proc/%u/ns/user", tpid) == -1) {
    perror("snprintf ns user path");
    return 1;
  }

  usrnsfd = open(usrnspath, O_RDONLY);
  if(usrnsfd == -1) {
    perror("open user namespace");
    return 1;
  }

  /* switch to container's user namespace so that user lookup returns correct uids */
  /* we allow this to fail if the container isn't user-namespaced */
  setns(usrnsfd, CLONE_NEWUSER);
  close(usrnsfd);

  /* switch to container's mount namespace/rootfs */
  if(setns(mntnsfd, CLONE_NEWNS) == -1) {
    perror("setns");
    return 1;
  }
  close(mntnsfd);

  printf("hello\n");

}
