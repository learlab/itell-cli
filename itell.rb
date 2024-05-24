# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Itell < Formula
  desc "Create and manage itell projects"
  homepage "https://github.com/learlab/itell-cli"
  version "0.0.1"

  on_macos do
    url "https://github.com/learlab/itell-cli/releases/download/v0.0.1/itell-cli_Darwin_all.tar.gz"
    sha256 "3d02b1d0b66b71366a3fa64032abd9cbfe188d0be1e241b78153da34624c0060"

    def install
      bin.install "itell-cli"
    end
  end

  on_linux do
    on_intel do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/learlab/itell-cli/releases/download/v0.0.1/itell-cli_Linux_x86_64.tar.gz"
        sha256 "bec33ed41d14f1ac8936ab9b580218ccb9405fb4e0377226d01ecc466a69c216"

        def install
          bin.install "itell-cli"
        end
      end
    end
    on_arm do
      if Hardware::CPU.is_64_bit?
        url "https://github.com/learlab/itell-cli/releases/download/v0.0.1/itell-cli_Linux_arm64.tar.gz"
        sha256 "d6276a551969ed82b369565081425c0d13630d627bdd81c5d8336f2fa0b37e14"

        def install
          bin.install "itell-cli"
        end
      end
    end
  end
end
